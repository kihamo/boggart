package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
	"github.com/kihamo/boggart/components/voice/players"
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logging.Logger
	provider    *yandex.YandexSpeechKitCloud
	audioPlayer *players.AudioPlayer
	routes      []dashboard.Route
}

func (c *Component) Name() string {
	return voice.ComponentName
}

func (c *Component) Version() string {
	return voice.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name:     mqtt.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.audioPlayer = players.NewAudio()
	c.logger = logging.DefaultLogger().Named(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.provider = yandex.NewYandexSpeechKitCloud(c.config.String(voice.ConfigYandexSpeechKitCloudKey), c.config.Bool(config.ConfigDebug))
	c.SetVolume(context.Background(), c.config.Int64(voice.ConfigSpeechVolume))

	<-a.ReadyComponent(mqtt.ComponentName)
	m := a.GetComponent(mqtt.ComponentName).(mqtt.Component)

	ready <- struct{}{}

	// audio updater
	var (
		lastStatus players.Status
		lastVolume int64
	)

	for {
		client := m.Client()

		if client != nil && client.IsConnected() {
			status := c.audioPlayer.Status()
			if status != lastStatus {
				m.Publish(voice.MQTTTopicPlayerStatus, 0, false, status.String())
				lastStatus = status
			}

			volume := c.audioPlayer.Volume()
			if volume != lastVolume {
				m.Publish(voice.MQTTTopicPlayerVolume, 0, false, fmt.Sprintf("%d", volume))
				lastVolume = volume
			}
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (c *Component) Speech(ctx context.Context, text string) error {
	return c.SpeechWithOptions(
		ctx,
		text,
		c.config.Int64(voice.ConfigSpeechVolume),
		c.config.Float64(voice.ConfigYandexSpeechKitCloudSpeed),
		c.config.String(voice.ConfigYandexSpeechKitCloudSpeaker))
}

func (c *Component) SpeechWithOptions(ctx context.Context, text string, volume int64, speed float64, speaker string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "speech_with_options")
	defer span.Finish()

	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}
	span.LogFields(log.Int64("volume", volume))

	if volume == 0 {
		c.logger.Warn("Skip speech text because volume is 0", "text", text)

		return nil
	}

	if speed < 0.1 {
		speed = 0.1
	} else if speed > 3 {
		speed = 3
	}
	span.LogFields(log.Float64("speed", speed))

	if speaker == "" {
		speaker = c.config.String(voice.ConfigYandexSpeechKitCloudSpeaker)
	}
	span.LogFields(log.String("speaker", speaker))

	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	c.logger.Debug("Speech text" + text)

	if c.provider == nil {
		err := errors.New("speech provider not found")

		tracing.SpanError(span, err)
		return err
	}

	file, err := c.provider.Generate(
		ctx,
		text,
		c.config.String(voice.ConfigYandexSpeechKitCloudLanguage),
		speaker,
		c.config.String(voice.ConfigYandexSpeechKitCloudEmotion),
		c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
		c.config.String(voice.ConfigYandexSpeechKitCloudQuality),
		speed)

	if err != nil {
		c.logger.Error("Error speech text",
			"text", text,
			"error", err.Error(),
		)

		tracing.SpanError(span, err)
		return err
	}

	err = c.SetVolume(ctx, volume)
	if err != nil {
		c.logger.Error("Failed set volume",
			"error", err.Error(),
			"format", c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text", text,
		)

		tracing.SpanError(span, err)
		return err
	}

	err = c.PlayReader(ctx, ioutil.NopCloser(bytes.NewReader(file)))
	if err != nil {
		c.logger.Error("Failed play speech text",
			"error", err.Error(),
			"format", c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text", text,
		)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}

func (c *Component) PlayReader(ctx context.Context, reader io.ReadCloser) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play.reader")
	defer span.Finish()

	err = c.audioPlayer.Stop()
	if err == nil {
		err = c.audioPlayer.PlayFromReader(reader)
	}

	if err != nil {
		c.logger.Error("Failed play reader", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play reader")
	}

	return err
}

func (c *Component) PlayURL(ctx context.Context, url string) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play.url")
	defer span.Finish()

	span.LogFields(log.String("url", url))

	err = c.audioPlayer.Stop()
	if err == nil {
		err = c.audioPlayer.PlayFromURL(url)
	}

	if err != nil {
		c.logger.Error("Failed play URL",
			"error", err.Error(),
			"url", url,
		)

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play URL", "url", url)
	}

	return err
}

func (c *Component) Play(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play")
	defer span.Finish()

	err := c.audioPlayer.Play()
	if err != nil {
		c.logger.Error("Failed play player", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play")
	}

	return err
}

func (c *Component) Pause(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.pause")
	defer span.Finish()

	err := c.audioPlayer.Pause()
	if err != nil {
		c.logger.Error("Failed pause player", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player pause")
	}

	return err
}

func (c *Component) Stop(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.stop")
	defer span.Finish()

	err := c.audioPlayer.Stop()
	if err != nil {
		c.logger.Error("Failed stop player", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player stopped")
	}

	return err
}

func (c *Component) Volume(ctx context.Context) int64 {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.volume.get")
	defer span.Finish()

	return c.audioPlayer.Volume()
}

func (c *Component) SetVolume(ctx context.Context, percent int64) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.volume.set")
	defer span.Finish()

	span.LogFields(log.Int64("percent", percent))

	err := c.audioPlayer.SetVolume(percent)
	if err != nil {
		c.logger.Error("Failed set volume player",
			"error", err.Error(),
			"volume", percent,
		)

		tracing.SpanError(span, err)
	} else {
		c.logger.Debugf("Player set volume %d", percent)
	}

	return err
}

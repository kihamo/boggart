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
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	mqtt        mqtt.Component
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
			Name: logger.ComponentName,
		},
		{
			Name:     mqtt.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.mqtt = a.GetComponent(mqtt.ComponentName).(mqtt.Component)
	c.audioPlayer = players.NewAudio()

	return nil
}

func (c *Component) Run() error {
	c.logger = logger.NewOrNop(c.Name(), c.application)
	c.provider = yandex.NewYandexSpeechKitCloud(c.config.String(voice.ConfigYandexSpeechKitCloudKey), c.config.Bool(config.ConfigDebug))
	c.application.GetComponent(mqtt.ComponentName).(mqtt.Component).Subscribe(NewMQTTSubscribe(c))
	c.audioPlayer.SetVolume(c.config.Int64(voice.ConfigSpeechVolume))

	c.audioPlayerUpdater()
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
		c.logger.Warn("Skip speech text because volume is 0", map[string]interface{}{
			"text": text,
		})

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

	c.logger.Debugf("Speech text %s", text)

	if c.provider == nil {
		err := errors.New("Speech provider not found")

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
		c.logger.Error("Error speech text", map[string]interface{}{
			"text":  text,
			"error": err.Error(),
		})

		return err
	}

	err = c.SetVolume(ctx, volume)
	if err != nil {
		c.logger.Error("Failed set volume", map[string]interface{}{
			"error":  err.Error(),
			"format": c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text":   text,
		})

		return err
	}

	err = c.PlayReader(ctx, ioutil.NopCloser(bytes.NewReader(file)))
	if err != nil {
		c.logger.Error("Failed play speech text", map[string]interface{}{
			"error":  err.Error(),
			"format": c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text":   text,
		})

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
		c.logger.Error("Failed play reader", map[string]interface{}{
			"error": err.Error(),
		})

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
		c.logger.Error("Failed play URL", map[string]interface{}{
			"error": err.Error(),
			"url":   url,
		})

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play URL", map[string]interface{}{
			"url": url,
		})
	}

	return err
}

func (c *Component) Play(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play")
	defer span.Finish()

	err := c.audioPlayer.Play()
	if err != nil {
		c.logger.Error("Failed play player", map[string]interface{}{
			"error": err.Error(),
		})

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
		c.logger.Error("Failed pause player", map[string]interface{}{
			"error": err.Error(),
		})

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
		c.logger.Error("Failed stop player", map[string]interface{}{
			"error": err.Error(),
		})

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
		c.logger.Error("Failed set volume player", map[string]interface{}{
			"error":  err.Error(),
			"volume": percent,
		})

		tracing.SpanError(span, err)
	} else {
		c.logger.Debugf("Player set volume %d", percent)
	}

	return err
}

func (c *Component) audioPlayerUpdater() {
	var (
		lastStatus players.Status
		lastVolume int64
	)

	for {
		client := c.mqtt.Client()

		if client != nil && client.IsConnected() {
			status := c.audioPlayer.Status()
			if status != lastStatus {
				c.mqtt.Publish(voice.MQTTTopicPlayerStatus, 0, false, status.String())
				lastStatus = status
			}

			volume := c.audioPlayer.Volume()
			if volume != lastVolume {
				c.mqtt.Publish(voice.MQTTTopicPlayerVolume, 0, false, fmt.Sprintf("%d", volume))
				lastVolume = volume
			}
		}

		time.Sleep(time.Second)
	}
}

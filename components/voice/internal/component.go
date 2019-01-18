package internal

import (
	"context"
	"errors"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/boggart/components/voice"
	"github.com/kihamo/boggart/components/voice/players"
	"github.com/kihamo/boggart/components/voice/players/alsa"
	"github.com/kihamo/boggart/components/voice/players/chromecast"
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"strconv"
	"strings"
	"time"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logging.Logger
	routes      []dashboard.Route

	textToSpeechProvider *yandex.YandexSpeechKitCloud
	players              map[string]players.Player
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
		{
			Name:     storage.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.players = make(map[string]players.Player, 0)
	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLogger().Named(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.textToSpeechProvider = yandex.NewYandexSpeechKitCloud(c.config.String(voice.ConfigYandexSpeechKitCloudKey)).
		WithDebug(c.config.Bool(config.ConfigDebug))

	if c.config.Bool(voice.ConfigPlayerALSAEnabled) {
		c.players["alsa"] = alsa.New()
	}

	addresses := c.config.String(voice.ConfigPlayerChromecastAddresses)
	if addresses != "" {
		for _, address := range strings.Split(addresses, ",") {
			address = strings.TrimSpace(address)
			if address == "" {
				c.logger.Warn("Chromecast address is empty")
				continue
			}

			parts := strings.SplitN(address, ":", 2)
			if len(parts) != 2 {
				c.logger.Warn("Bad Chromecast address " + address)
				continue
			}

			playerID := mqtt.NameReplace("chromecast/" + strings.ToLower(parts[0]))
			port, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				c.logger.Warn("Bad Chromecast address "+address, "error", err.Error())
				continue
			}

			c.players[playerID] = chromecast.New(parts[0], port)
		}
	}

	if len(c.players) > 0 {
		go c.playersUpdater()
	}

	return nil
}

func (c *Component) playersUpdater() {
	<-c.application.ReadyComponent(mqtt.ComponentName)
	m := c.application.GetComponent(mqtt.ComponentName).(mqtt.Component)

	storeLastStatus := make(map[string]int64, len(c.players))
	storeLastVolume := make(map[string]int64, len(c.players))
	storeLastMute := make(map[string]bool, len(c.players))

	for {
		client := m.Client()
		if client != nil && client.IsConnected() {
			for name, player := range c.players {
				status := player.Status()
				lastStatus, ok := storeLastStatus[name]

				if !ok || status.Int64() != lastStatus {
					err := m.Publish(context.Background(), MQTTTopicPlayerStateStatus.Format(name), 0, false, status.String())
					if err == nil {
						storeLastStatus[name] = status.Int64()
					}
				}

				volume, err := player.Volume()
				if err == nil {
					lastVolume, ok := storeLastVolume[name]

					if !ok || volume != lastVolume {
						err := m.Publish(context.Background(), MQTTTopicPlayerStateVolume.Format(name), 0, false, strconv.FormatInt(volume, 10))
						if err == nil {
							storeLastVolume[name] = volume
						}
					}
				}

				mute, err := player.Mute()
				if err == nil {
					lastMute, ok := storeLastMute[name]

					if !ok || mute != lastMute {
						err := m.Publish(context.Background(), MQTTTopicPlayerStateMute.Format(name), 0, false, mute)
						if err == nil {
							storeLastMute[name] = mute
						}
					}
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func (c *Component) Players() map[string]players.Player {
	return c.players
}

func (c *Component) Speech(ctx context.Context, player string, text string) error {
	return c.SpeechWithOptions(
		ctx,
		player,
		text,
		c.config.Int64(voice.ConfigSpeechVolume),
		c.config.Float64(voice.ConfigYandexSpeechKitCloudSpeed),
		c.config.String(voice.ConfigYandexSpeechKitCloudSpeaker))
}

func (c *Component) SpeechWithOptions(ctx context.Context, player string, text string, volume int64, speed float64, speaker string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "speech_with_options")
	defer span.Finish()

	span.LogFields(log.String("player", player))

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

	if c.textToSpeechProvider == nil {
		err := errors.New("text to speech provider not found")

		tracing.SpanError(span, err)
		return err
	}

	err := c.SetVolume(ctx, player, volume)
	if err != nil {
		c.logger.Error("Failed set volume",
			"error", err.Error(),
			"format", c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text", text,
			"player", player,
		)

		tracing.SpanError(span, err)
		return err
	}

	u := c.textToSpeechProvider.GenerateURL(
		ctx,
		text,
		c.config.String(voice.ConfigYandexSpeechKitCloudLanguage),
		speaker,
		c.config.String(voice.ConfigYandexSpeechKitCloudEmotion),
		c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
		c.config.String(voice.ConfigYandexSpeechKitCloudQuality),
		speed)

	err = c.PlayURL(ctx, player, u)
	if err != nil {
		c.logger.Error("Failed play speech text",
			"error", err.Error(),
			"format", c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text", text,
			"player", player,
		)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}

func (c *Component) PlayReader(ctx context.Context, player string, reader io.ReadCloser) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play.reader")
	defer span.Finish()

	if p, ok := c.players[player]; ok {
		err = p.Stop()
		if err == nil {
			err = p.PlayFromReader(reader)
		}
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed play reader", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play reader")
	}

	return err
}

func (c *Component) PlayURL(ctx context.Context, player string, url string) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play.url")
	defer span.Finish()

	span.LogFields(log.String("url", url))

	if p, ok := c.players[player]; ok {
		err = p.Stop()
		if err == nil {
			err = p.PlayFromURL(url)
		}
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed play URL",
			"error", err.Error(),
			"url", url,
			"player", player,
		)

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play URL", "url", url)
	}

	return err
}

func (c *Component) Play(ctx context.Context, player string) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.play")
	defer span.Finish()

	if p, ok := c.players[player]; ok {
		err = p.Play()
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed play player", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player play")
	}

	return err
}

func (c *Component) Pause(ctx context.Context, player string) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.pause")
	defer span.Finish()

	if p, ok := c.players[player]; ok {
		err = p.Pause()
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed pause player", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player pause")
	}

	return err
}

func (c *Component) Stop(ctx context.Context, player string) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.stop")
	defer span.Finish()

	if p, ok := c.players[player]; ok {
		err = p.Stop()
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed stop player", "error", err.Error())

		tracing.SpanError(span, err)
	} else {
		c.logger.Debug("Player stopped")
	}

	return err
}

func (c *Component) Volume(ctx context.Context, player string) (volume int64, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.volume.get")
	defer span.Finish()

	if p, ok := c.players[player]; ok {
		volume, err = p.Volume()
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed get player volume", "error", err.Error())

		tracing.SpanError(span, err)
	}

	return volume, err
}

func (c *Component) SetVolume(ctx context.Context, player string, percent int64) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.volume.set")
	defer span.Finish()

	span.LogFields(log.Int64("percent", percent))

	if p, ok := c.players[player]; ok {
		err = p.SetVolume(percent)
	} else {
		err = errors.New("Player " + player + "not found")
	}

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

func (c *Component) SetMute(ctx context.Context, player string, mute bool) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, voice.ComponentName, "player.mute.set")
	defer span.Finish()

	span.LogFields(log.Bool("mute", mute))

	if p, ok := c.players[player]; ok {
		err = p.SetMute(mute)
	} else {
		err = errors.New("Player " + player + "not found")
	}

	if err != nil {
		c.logger.Error("Failed set mute player",
			"error", err.Error(),
			"mute", mute,
		)

		tracing.SpanError(span, err)
	} else {
		c.logger.Debugf("Player set mute %v", mute)
	}

	return err
}

func (c *Component) TextToSpeechProvider() *yandex.YandexSpeechKitCloud {
	return c.textToSpeechProvider
}

func (c *Component) Shutdown() error {
	for _, player := range c.players {
		player.Close()
	}

	return nil
}

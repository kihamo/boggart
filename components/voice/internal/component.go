package internal

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	provider    *yandex.YandexSpeechKitCloud
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

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.Name(), c.application)
	c.provider = yandex.NewYandexSpeechKitCloud(c.config.String(voice.ConfigYandexSpeechKitCloudKey), c.config.Bool(config.ConfigDebug))
	c.application.GetComponent(mqtt.ComponentName).(mqtt.Component).Subscribe(NewMQTTSubscribe(c))

	return nil
}

func (c *Component) Speech(text string) error {
	volumePercent := c.config.Int64(voice.ConfigSpeechVolume)
	if volumePercent < 0 {
		volumePercent = 0
	} else if volumePercent > 100 {
		volumePercent = 100
	}

	if volumePercent == 0 {
		c.logger.Error("Skip speech text because volume is 0", map[string]interface{}{
			"text": text,
		})

		return nil
	}

	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	c.logger.Debugf("Speech text %s", text)

	if c.provider == nil {
		return errors.New("Speech provider not found")
	}

	file, err := c.provider.Generate(
		context.Background(),
		text,
		c.config.String(voice.ConfigYandexSpeechKitCloudLanguage),
		c.config.String(voice.ConfigYandexSpeechKitCloudSpeaker),
		c.config.String(voice.ConfigYandexSpeechKitCloudEmotion),
		c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
		c.config.String(voice.ConfigYandexSpeechKitCloudQuality),
		c.config.Float64(voice.ConfigYandexSpeechKitCloudSpeed))

	if err != nil {
		c.logger.Error("Error speech text", map[string]interface{}{
			"text": text,
		})

		return err
	}

	f := ioutil.NopCloser(bytes.NewReader(file))

	// decode
	var (
		stream beep.StreamSeekCloser
		format beep.Format
	)

	switch c.config.String(voice.ConfigYandexSpeechKitCloudFormat) {
	case yandex.FormatMP3:
		stream, format, err = mp3.Decode(f)

	case yandex.FormatWAV:
		stream, format, err = wav.Decode(f)

	default:
		err = errors.New("Unknown format of audio file")
	}

	f.Close()

	if err != nil {
		c.logger.Error("Failed decode audio file for speech", map[string]interface{}{
			"error":  err.Error(),
			"format": c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
			"text":   text,
		})

		return err
	}

	// sound effects
	streamWithEffects := effects.Volume{
		Streamer: stream,
		Base:     2,
		Volume:   -float64(100-volumePercent) / 100.0 * 5,
		Silent:   false,
	}

	// play
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(&streamWithEffects)

	return nil
}

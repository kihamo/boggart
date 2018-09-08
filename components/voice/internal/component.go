package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
	"github.com/kihamo/boggart/components/voice/players"
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	mqtt        mqtt.Component
	provider    *yandex.YandexSpeechKitCloud
	audioPlayer *players.AudioPlayer
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

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.Name(), c.application)
	c.provider = yandex.NewYandexSpeechKitCloud(c.config.String(voice.ConfigYandexSpeechKitCloudKey), c.config.Bool(config.ConfigDebug))
	c.application.GetComponent(mqtt.ComponentName).(mqtt.Component).Subscribe(NewMQTTSubscribe(c))
	c.audioPlayer.SetVolume(c.config.Int64(voice.ConfigSpeechVolume))

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.audioPlayerUpdater()
	}()

	return nil
}

func (c *Component) Speech(text string) error {
	return c.SpeechWithOptions(
		text,
		c.config.Int64(voice.ConfigSpeechVolume),
		c.config.Float64(voice.ConfigYandexSpeechKitCloudSpeed),
		c.config.String(voice.ConfigYandexSpeechKitCloudSpeaker))
}

func (c *Component) SpeechWithOptions(text string, volume int64, speed float64, speaker string) error {
	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}

	if volume == 0 {
		c.logger.Error("Skip speech text because volume is 0", map[string]interface{}{
			"text": text,
		})

		return nil
	}

	if speed < 0.1 {
		speed = 0.1
	} else if speed > 3 {
		speed = 3
	}

	if speaker == "" {
		speaker = c.config.String(voice.ConfigYandexSpeechKitCloudSpeaker)
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
		speaker,
		c.config.String(voice.ConfigYandexSpeechKitCloudEmotion),
		c.config.String(voice.ConfigYandexSpeechKitCloudFormat),
		c.config.String(voice.ConfigYandexSpeechKitCloudQuality),
		speed)

	if err != nil {
		c.logger.Error("Error speech text", map[string]interface{}{
			"text": text,
		})

		return err
	}

	c.audioPlayer.SetVolume(volume)

	err = c.PlayReader(ioutil.NopCloser(bytes.NewReader(file)))
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

func (c *Component) PlayReader(reader io.ReadCloser) (err error) {
	err = c.audioPlayer.Stop()
	if err == nil {
		err = c.audioPlayer.PlayFromReader(reader)
	}

	if err != nil {
		c.logger.Error("Failed play reader", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		c.logger.Debug("Player play reader")
	}

	return err
}

func (c *Component) PlayURL(url string) (err error) {
	err = c.audioPlayer.Stop()
	if err == nil {
		err = c.audioPlayer.PlayFromURL(url)
	}

	if err != nil {
		c.logger.Error("Failed play URL", map[string]interface{}{
			"error": err.Error(),
			"url":   url,
		})
	} else {
		c.logger.Debug("Player play URL", map[string]interface{}{
			"url": url,
		})
	}

	return err
}

func (c *Component) Play() error {
	err := c.audioPlayer.Play()
	if err != nil {
		c.logger.Error("Failed play player", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		c.logger.Debug("Player play")
	}

	return err
}

func (c *Component) Pause() error {
	err := c.audioPlayer.Pause()
	if err != nil {
		c.logger.Error("Failed pause player", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		c.logger.Debug("Player pause")
	}

	return err
}

func (c *Component) Stop() error {
	err := c.audioPlayer.Stop()
	if err != nil {
		c.logger.Error("Failed stop player", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		c.logger.Debug("Player stopped")
	}

	return err
}

func (c *Component) Volume() int64 {
	return c.audioPlayer.Volume()
}

func (c *Component) SetVolume(percent int64) error {
	err := c.audioPlayer.SetVolume(percent)
	if err != nil {
		c.logger.Error("Failed set volume player", map[string]interface{}{
			"error":  err.Error(),
			"volume": percent,
		})
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

		time.Sleep(time.Second)
	}
}

package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"

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
	c.logger.Debugf("Speach text %s", text)

	if c.provider == nil {
		return errors.New("Provider not found")
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
		return err
	}

	fmt.Println(file)

	return nil
}

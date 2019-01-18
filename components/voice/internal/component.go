package internal

import (
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/boggart/components/voice"
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logging.Logger
	routes      []dashboard.Route

	textToSpeechProvider *yandex.YandexSpeechKitCloud
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
			Name:     storage.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	c.logger = logging.DefaultLogger().Named(c.Name())

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.textToSpeechProvider = yandex.NewYandexSpeechKitCloud(c.config.String(voice.ConfigYandexSpeechKitCloudKey)).
		WithDebug(c.config.Bool(config.ConfigDebug))

	return nil
}

func (c *Component) TextToSpeechProvider() *yandex.YandexSpeechKitCloud {
	return c.textToSpeechProvider
}

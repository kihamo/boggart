package internal

import (
	"io/ioutil"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	_ "github.com/kihamo/boggart/components/boggart/bind/broadlink"
	_ "github.com/kihamo/boggart/components/boggart/bind/ds18b20"
	_ "github.com/kihamo/boggart/components/boggart/bind/google_home"
	_ "github.com/kihamo/boggart/components/boggart/bind/gpio"
	_ "github.com/kihamo/boggart/components/boggart/bind/hikvision"
	_ "github.com/kihamo/boggart/components/boggart/bind/led_wifi"
	_ "github.com/kihamo/boggart/components/boggart/bind/lg_webos"
	_ "github.com/kihamo/boggart/components/boggart/bind/mercury"
	_ "github.com/kihamo/boggart/components/boggart/bind/mikrotik"
	_ "github.com/kihamo/boggart/components/boggart/bind/nut"
	_ "github.com/kihamo/boggart/components/boggart/bind/pulsar"
	_ "github.com/kihamo/boggart/components/boggart/bind/samsung_tizen"
	_ "github.com/kihamo/boggart/components/boggart/bind/softvideo"
	"github.com/kihamo/boggart/components/boggart/internal/manager"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/syslog"
	w "github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
	"gopkg.in/yaml.v2"
	"periph.io/x/periph/host"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logging.Logger
	routes      []dashboard.Route

	listenersManager *w.ListenersManager
	manager          *manager.Manager
}

type FileYAML struct {
	Devices []BindYAML
}

type BindYAML struct {
	Enabled     *bool
	Type        string
	ID          *string
	Description string
	Tags        []string
	Config      map[string]interface{}
}

func (c *Component) Name() string {
	return boggart.ComponentName
}

func (c *Component) Version() string {
	return boggart.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name: annotations.ComponentName,
		},
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logging.ComponentName,
		},
		{
			Name: messengers.ComponentName,
		},
		{
			Name:     mqtt.ComponentName,
			Required: true,
		},
		{
			Name:     metrics.ComponentName,
			Required: true,
		},
		{
			Name: syslog.ComponentName,
		},
		{
			Name:     workers.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.listenersManager = w.NewListenersManager()

	return nil
}

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	<-a.ReadyComponent(mqtt.ComponentName)
	<-a.ReadyComponent(workers.ComponentName)

	c.manager = manager.NewManager(
		a.GetComponent(mqtt.ComponentName).(mqtt.Component),
		a.GetComponent(workers.ComponentName).(workers.Component),
		c.listenersManager)

	c.logger = logging.DefaultLogger().Named(c.Name())

	if _, err := host.Init(); err != nil {
		return err
	}

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	c.ReloadConfig()
	c.manager.Ready()

	return nil
}

func (c *Component) ReloadConfig() error {
	if err := c.manager.UnregisterAll(); err != nil {
		return err
	}

	if err := c.initConfigFromYaml(); err != nil {
		return err
	}

	return nil
}

func (c *Component) initConfigFromYaml() error {
	fileName := c.config.String(boggart.ConfigConfigYAML)
	if fileName == "" {
		return nil
	}

	var fileYAML FileYAML

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &fileYAML)
	if err != nil {
		return err
	}

	for _, d := range fileYAML.Devices {
		if d.Enabled != nil && !*d.Enabled {
			continue
		}

		if d.Type == "" {
			// TODO: error
			continue
		}

		kind, err := boggart.GetBindType(d.Type)
		if err != nil {
			return err
		}

		cfg, err := boggart.ValidateBindConfig(kind, d.Config)
		if err != nil {
			return err
		}

		bind, err := kind.CreateBind(cfg)
		if err != nil {
			return err
		}

		if d.ID != nil && *d.ID != "" {
			_, err = c.manager.RegisterWithID(*d.ID, bind, d.Type, d.Description, d.Tags, cfg)
		} else {
			_, err = c.manager.Register(bind, d.Type, d.Description, d.Tags, cfg)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

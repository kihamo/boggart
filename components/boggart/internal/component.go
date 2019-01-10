package internal

import (
	"io/ioutil"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/syslog"
	"github.com/kihamo/go-workers/manager"
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

	listenersManager *manager.ListenersManager
	devicesManager   *DevicesManager
}

type FileYAML struct {
	Devices []DeviceYAML
}

type DeviceYAML struct {
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
	return nil
}

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	c.listenersManager = manager.NewListenersManager()

	<-a.ReadyComponent(workers.ComponentName)
	c.devicesManager = NewDevicesManager(
		a.GetComponent(mqtt.ComponentName).(mqtt.Component),
		a.GetComponent(workers.ComponentName).(workers.Component),
		c.listenersManager)

	c.logger = logging.DefaultLogger().Named(c.Name())

	if _, err := host.Init(); err != nil {
		return err
	}

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	err := c.initConfigFromYaml()
	if err != nil {
		return err
	}

	c.devicesManager.Ready()

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

		kind, err := boggart.GetDeviceType(d.Type)
		if err != nil {
			return err
		}

		device, err := kind.CreateBind(d.Config)
		if err != nil {
			return err
		}

		if d.ID != nil && *d.ID != "" {
			c.devicesManager.RegisterWithID(*d.ID, device, d.Type, d.Description, d.Tags, d.Config)
		} else {
			c.devicesManager.Register(device, d.Type, d.Description, d.Tags, d.Config)
		}
	}

	return nil
}

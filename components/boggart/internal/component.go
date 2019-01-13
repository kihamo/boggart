package internal

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kihamo/boggart/components/boggart"
	_ "github.com/kihamo/boggart/components/boggart/bind/broadlink"
	_ "github.com/kihamo/boggart/components/boggart/bind/ds18b20"
	_ "github.com/kihamo/boggart/components/boggart/bind/google_home_mini"
	_ "github.com/kihamo/boggart/components/boggart/bind/gpio"
	_ "github.com/kihamo/boggart/components/boggart/bind/hikvision"
	_ "github.com/kihamo/boggart/components/boggart/bind/led_wifi"
	_ "github.com/kihamo/boggart/components/boggart/bind/lg_webos"
	_ "github.com/kihamo/boggart/components/boggart/bind/mercury"
	_ "github.com/kihamo/boggart/components/boggart/bind/mikrotik"
	_ "github.com/kihamo/boggart/components/boggart/bind/nut"
	_ "github.com/kihamo/boggart/components/boggart/bind/pulsar_heat_meter"
	_ "github.com/kihamo/boggart/components/boggart/bind/samsung_tizen"
	_ "github.com/kihamo/boggart/components/boggart/bind/softvideo"
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
	"github.com/mitchellh/mapstructure"
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
	c.listenersManager = manager.NewListenersManager()

	return nil
}

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	<-a.ReadyComponent(mqtt.ComponentName)
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

		var cfg interface{}

		if prepare := kind.Config(); prepare != nil {
			mapStructureDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   &prepare,
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeHookFunc(time.RFC3339),
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToIPHookFunc(),
					mapstructure.StringToIPNetHookFunc(),
				),
			})

			if err != nil {
				return err
			}

			if err := mapStructureDecoder.Decode(d.Config); err != nil {
				return err
			}

			if _, err = govalidator.ValidateStruct(prepare); err != nil {
				return err
			}

			cfg = prepare
		} else {
			cfg = d.Config
		}

		device, err := kind.CreateBind(cfg)
		if err != nil {
			return err
		}

		if d.ID != nil && *d.ID != "" {
			err = c.devicesManager.RegisterWithID(*d.ID, device, d.Type, d.Description, d.Tags, cfg)
		} else {
			_, err = c.devicesManager.Register(device, d.Type, d.Description, d.Tags, cfg)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

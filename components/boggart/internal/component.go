package internal

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
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
	"periph.io/x/periph/host"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logging.Logger
	routes      []dashboard.Route

	connectionRS485  *rs485.Connection
	listenersManager *manager.ListenersManager
	devicesManager   *DevicesManager
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

	c.devicesManager.SetCheckerTickerDuration(c.config.Duration(boggart.ConfigDevicesManagerCheckInterval))
	c.devicesManager.SetCheckerTimeout(c.config.Duration(boggart.ConfigDevicesManagerCheckTimeout))

	c.initListeners()
	c.initRS485()

	c.initGPIO()
	c.initElectricityMeters()
	c.initInternetProviders()
	c.initPhones()
	c.initRouters()
	c.initCameras()
	c.initPulsarMeters()
	c.initSensor()
	c.initSockets()
	c.initRemoteControl()

	c.devicesManager.Ready()

	return nil
}

func (c *Component) initRS485() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connectionRS485 = rs485.NewConnection(
		c.config.String(boggart.ConfigRS485Address),
		c.config.Duration(boggart.ConfigRS485Timeout))
}

func (c *Component) RS485() *rs485.Connection {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.connectionRS485
}

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
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	config      config.Component
	logger      logger.Logger
	workers     workers.Component
	routes      []dashboard.Route

	connectionRS485  *rs485.Connection
	listenersManager *manager.ListenersManager
	devicesManager   *DevicesManager
	// securityManager  *SecurityManager
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
			Name: logger.ComponentName,
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

	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.workers = a.GetComponent(workers.ComponentName).(workers.Component)

	c.listenersManager = manager.NewListenersManager()
	c.devicesManager = NewDevicesManager(c.workers, c.listenersManager)
	// c.securityManager = NewSecurityManager(c.devicesManager, c.listenersManager)

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.Name(), c.application)
	c.devicesManager.SetCheckerTickerDuration(c.config.Duration(boggart.ConfigDevicesManagerCheckInterval))
	c.devicesManager.SetCheckerTimeout(c.config.Duration(boggart.ConfigDevicesManagerCheckTimeout))

	c.initListeners()
	c.initConnectionRS485()

	//c.initPC()
	c.initGPIO()
	//c.initCameras()
	c.initElectricityMeters()
	c.initInternetProviders()
	c.initPhones()
	c.initRouters()
	c.initVideoRecorders()
	c.initPulsarMeters()
	//c.initUPS()
	c.initTV()
	// c.initLight()
	c.initSensor()

	c.initMQTT()

	c.devicesManager.Ready()

	return nil
}

func (c *Component) initConnectionRS485() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connectionRS485 = rs485.NewConnection(
		c.config.String(boggart.ConfigRS485Address),
		c.config.Duration(boggart.ConfigRS485Timeout))
}

func (c *Component) ConnectionRS485() *rs485.Connection {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.connectionRS485
}

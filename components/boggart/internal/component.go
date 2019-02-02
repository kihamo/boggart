package internal

import (
	"errors"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/manager"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/syslog"
	w "github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
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
	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	<-a.ReadyComponent(mqtt.ComponentName, workers.ComponentName)

	var i18nCmp i18n.Component
	if a.HasComponent(i18n.ComponentName) {
		<-a.ReadyComponent(i18n.ComponentName)
		i18nCmp = a.GetComponent(i18n.ComponentName).(i18n.Component)
	}

	c.mutex.Lock()
	c.manager = manager.NewManager(
		a.GetComponent(dashboard.ComponentName).(dashboard.Component),
		i18nCmp,
		a.GetComponent(mqtt.ComponentName).(mqtt.Component),
		a.GetComponent(workers.ComponentName).(workers.Component),
		c.listenersManager)
	c.mutex.Unlock()

	c.logger = logging.DefaultLogger().Named(c.Name())

	if _, err := host.Init(); err != nil {
		return err
	}

	if loaded, err := c.ReloadConfig(); err == nil {
		c.logger.Debug("Loaded " + strconv.FormatInt(int64(loaded), 10) + " binds from " + c.config.String(boggart.ConfigConfigYAML) + " file")
	} else {
		return err
	}

	c.manager.Ready()

	return nil
}

func (c *Component) ReloadConfig() (int, error) {
	if err := c.manager.UnregisterAll(); err != nil {
		return -1, err
	}

	loaded, err := c.initConfigFromYaml("")
	if err != nil {
		return -1, err
	}

	return loaded, nil
}

func (c *Component) ReloadConfigByID(id string) error {
	if err := c.manager.Unregister(id); err != nil {
		return err
	}

	loaded, err := c.initConfigFromYaml(id)
	if err != nil {
		return err
	}

	if loaded == 0 {
		return errors.New("id " + id + " not found in file")
	}

	return nil
}

func (c *Component) initConfigFromYaml(id string) (int, error) {
	fileName := c.config.String(boggart.ConfigConfigYAML)
	if fileName == "" {
		return -1, nil
	}

	var fileYAML FileYAML

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, err
	}

	err = yaml.Unmarshal(data, &fileYAML)
	if err != nil {
		return -1, err
	}

	var loaded int
	for _, d := range fileYAML.Devices {
		if id != "" {
			if d.ID == nil || *d.ID == "" || *d.ID != id {
				continue
			}
		}

		if d.Enabled != nil && !*d.Enabled {
			continue
		}

		if d.Type == "" {
			return -1, errors.New("empty type")
		}

		kind, err := boggart.GetBindType(d.Type)
		if err != nil {
			return -1, err
		}

		cfg, err := boggart.ValidateBindConfig(kind, d.Config)
		if err != nil {
			return -1, err
		}

		bind, err := kind.CreateBind(cfg)
		if err != nil {
			return -1, err
		}

		if d.ID != nil && *d.ID != "" {
			_, err = c.manager.RegisterWithID(*d.ID, bind, d.Type, d.Description, d.Tags, cfg)
		} else {
			_, err = c.manager.Register(bind, d.Type, d.Description, d.Tags, cfg)
		}

		if err != nil {
			return -1, err
		}

		loaded++
	}

	return loaded, nil
}

func (c *Component) Shutdown() error {
	c.mutex.Lock()
	m := c.manager
	c.mutex.Unlock()

	if m != nil {
		return m.UnregisterAll()
	}

	return nil
}

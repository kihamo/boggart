package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	bb "github.com/kihamo/boggart/components/boggart/bind/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/pborman/uuid"
	"gopkg.in/yaml.v2"
	"periph.io/x/periph/host"
)

type Component struct {
	application shadow.Application
	config      config.Component
	dashboard   dashboard.Component
	i18n        i18n.Component
	logger      logging.Logger
	mqtt        mqtt.Component

	routes       []dashboard.Route
	binds        sync.Map
	closing      int32
	tasksManager *tasks.Manager
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
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.tasksManager = tasks.NewManager()

	return nil
}

func (c *Component) Run(a shadow.Application, _ chan<- struct{}) error {
	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	<-a.ReadyComponent(mqtt.ComponentName)
	c.mqtt = a.GetComponent(mqtt.ComponentName).(mqtt.Component)
	c.mqtt.OnConnectHandlerAdd(c.mqttOnConnectHandler)

	if a.HasComponent(i18n.ComponentName) {
		<-a.ReadyComponent(i18n.ComponentName)
		c.i18n = a.GetComponent(i18n.ComponentName).(i18n.Component)
	}

	c.dashboard = a.GetComponent(dashboard.ComponentName).(dashboard.Component)
	c.logger = logging.DefaultLazyLogger(c.Name())

	if _, err := host.Init(); err != nil {
		return err
	}

	if loaded, err := c.ReloadConfig(); err == nil {
		c.logger.Debug("Loaded " + strconv.FormatInt(int64(loaded), 10) + " binds from " + c.config.String(boggart.ConfigConfigYAML) + " file")
	} else {
		return err
	}

	return nil
}

func (c *Component) registerDefaultBinds() (int, error) {
	kind, err := boggart.GetBindType(c.Name())
	if err != nil {
		return -1, err
	}

	cfg, _, err := boggart.ValidateBindConfig(kind, map[string]interface{}{
		"application_name":    c.application.Name(),
		"application_version": c.application.Version(),
		"application_build":   c.application.Build(),
	})
	if err != nil {
		return -1, err
	}

	bind, err := kind.CreateBind(cfg)
	if err != nil {
		return -1, err
	}

	if bindSupport, ok := bind.(*bb.Bind); ok {
		bindSupport.SetApplication(c.application)
	}

	bindID := c.config.String(boggart.ConfigBoggartBindID)

	_, err = c.RegisterBind(bindID, bind, c.Name(), c.application.Name(), []string{c.Name()}, cfg)
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (c *Component) ReloadConfig() (int, error) {
	if err := c.unregisterBindsAll(); err != nil {
		return -1, err
	}

	loadedDefault, err := c.registerDefaultBinds()
	if err != nil {
		return -1, err
	}

	loadedFromConfig, err := c.initConfigFromYaml("")
	if err != nil {
		return -1, err
	}

	return loadedDefault + loadedFromConfig, nil
}

func (c *Component) ReloadConfigByID(id string) error {
	if err := c.UnregisterBindByID(id); err != nil {
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

		cfg, md, err := boggart.ValidateBindConfig(kind, d.Config)
		if err != nil {
			return -1, fmt.Errorf("config of device type %s validate failed with error: %v", d.Type, err)
		}

		bind, err := kind.CreateBind(cfg)
		if err != nil {
			return -1, err
		}

		var id string
		if d.ID != nil && *d.ID != "" {
			id = *d.ID
		}

		if _, err := c.RegisterBind(id, bind, d.Type, d.Description, d.Tags, cfg); err != nil {
			return -1, err
		}

		if len(md.Unused) > 0 {
			if logger, ok := di.LoggerContainerBind(bind); ok {
				for _, field := range md.Unused {
					logger.Warn("Unused config field", "field", field)
				}
			}
		}

		loaded++
	}

	return loaded, nil
}

func (c *Component) Shutdown() error {
	atomic.StoreInt32(&c.closing, 1)

	return c.unregisterBindsAll()
}

func (c *Component) RegisterBind(id string, bind boggart.Bind, t string, description string, tags []string, config interface{}) (boggart.BindItem, error) {
	if id == "" {
		id = uuid.New()
	} else if existsBind := c.Bind(id); existsBind != nil {
		return nil, errors.New("bind item with id " + id + " already exist")
	}

	bindType, err := boggart.GetBindType(t)
	if err != nil {
		return nil, err
	}

	bindItem := &BindItem{
		bind:        bind,
		bindType:    bindType,
		id:          id,
		t:           t,
		description: description,
		tags:        tags,
		config:      config,
	}

	c.itemStatusUpdate(bindItem, boggart.BindStatusInitializing)

	c.binds.Store(id, bindItem)

	c.logger.Debug("Register bind", "type", bindItem.Type(), "id", bindItem.ID())

	go func() {
		var err error

		defer func() {
			if err != nil {
				c.itemStatusUpdate(bindItem, boggart.BindStatusUninitialized)
				c.itemLogger(bind).Error("Bind run failed",
					"type", bindItem.Type(),
					"id", bindItem.ID(),
					"err", err.Error(),
				)
			}
		}()

		// mqtt container
		if bindSupport, ok := bind.(di.MQTTContainerSupport); ok {
			bindSupport.SetMQTT(di.NewMQTTContainer(bindItem, c.mqtt))
		}

		// config container
		if bindSupport, ok := bind.(di.ConfigContainerSupport); ok {
			bindSupport.SetConfig(di.NewConfigContainer(bindItem, c.config))
		}

		// meta container
		if bindSupport, ok := bind.(di.MetaContainerSupport); ok {
			bindSupport.SetMeta(di.NewMetaContainer(bindItem, c.mqtt, c.config))
		}

		// logger container
		if bindSupport, ok := bind.(di.LoggerContainerSupport); ok {
			bindSupport.SetLogger(di.NewLoggerContainer(bindItem, c.logger))
		}

		// metrics container
		if bindSupport, ok := bind.(di.MetricsContainerSupport); ok {
			bindSupport.SetMetrics(di.NewMetricsContainer(bindItem))
		}

		// widget container
		if bindSupport, ok := bind.(di.WidgetContainerSupport); ok {
			ctr := di.NewWidgetContainer(bindItem, c.config)

			bindSupport.SetWidget(ctr)

			if fs := ctr.AssetFS(); fs != nil {
				name := ctr.TemplateNamespace()

				// templates
				render := c.dashboard.Renderer()
				if !render.IsRegisterNamespace(name) {
					if err = render.RegisterNamespace(name, fs); err != nil {
						err = fmt.Errorf("bind templates failed: %w", err)
						return
					}
				}

				// asset fs
				c.dashboard.RegisterAssetFS(name, fs)

				// i18n
				if c.i18n != nil {
					fs.Prefix = "locales"
					c.i18n.LoadLocaleFromFiles(name, i18n.FromAssetFS(fs))
				}
			}
		}

		if runner, ok := bind.(boggart.BindRunnable); ok {
			if err = runner.Run(); err != nil {
				return
			}
		}

		// mqtt subscribers
		if bindSupport, ok := di.MQTTContainerBind(bind); ok {
			// TODO: обвешать подписки враппером, что бы только в online можно было посылать
			for _, subscriber := range bindSupport.Subscribers() {
				if err = c.mqtt.SubscribeSubscriber(subscriber.Subscriber()); err != nil {
					err = fmt.Errorf("bind mqtt subscribe %s failed: %w",
						subscriber.Subscriber().Topic().String(),
						err,
					)
					return
				}
			}
		}

		// workers container
		if bindSupport, ok := bindItem.Bind().(di.WorkersContainerSupport); ok {
			ctr := di.NewWorkersContainer(bindItem, c.tasksManager)

			bindSupport.SetWorkers(ctr)

			if err = ctr.HookRegister(); err != nil {
				err = fmt.Errorf("register tasks failed: %w", err)
				return
			}
		}

		// probes container
		if bindSupport, ok := bind.(di.ProbesContainerSupport); ok {
			ctr := di.NewProbesContainer(
				bindItem,
				func(status boggart.BindStatus) {
					c.itemStatusUpdate(bindItem, status)
				}, func() error {
					_, err := c.RegisterBind(bindItem.ID(), bindItem.Bind(), bindItem.Type(), bindItem.Description(), bindItem.Tags(), bindItem.Config())
					return err
				}, func() error {
					return c.UnregisterBindByID(bindItem.ID())
				},
				c.tasksManager,
				metricProbes)

			bindSupport.SetProbes(ctr)

			if err = ctr.HookRegister(); err != nil {
				err = fmt.Errorf("register probes failed: %w", err)
				return
			}
		} else {
			c.itemStatusUpdate(bindItem, boggart.BindStatusOnline)
		}
	}()

	return bindItem, nil
}

func (c *Component) UnregisterBindByID(id string) error {
	d, ok := c.binds.Load(id)
	if !ok {
		return nil
	}

	bindItem := d.(*BindItem)
	bind := bindItem.Bind()
	logger := c.itemLogger(bind)

	c.itemStatusUpdate(bindItem, boggart.BindStatusRemoving)

	// unregister mqtt
	if bindSupport, ok := di.MQTTContainerBind(bind); ok {
		// не блокирующее отписываемся, так как mqtt может быть не доступен
		go func() {
			for _, subscriber := range bindSupport.Subscribers() {
				if err := c.mqtt.UnsubscribeSubscriber(subscriber.Subscriber()); err != nil {
					logger.Error("Unregister bind failed because unsubscribe MQTT failed",
						"type", bindItem.Type(),
						"id", bindItem.ID(),
						"error", err.Error(),
					)
				}
			}

			if s := bindItem.Status(); s.IsStatusRemoving() || s.IsStatusRemoved() {
				bindSupport.SetClient(nil)
			}
		}()
	}

	// remove probes
	if bindSupport, ok := di.ProbesContainerBind(bind); ok {
		bindSupport.HookUnregister()
	}

	// workers
	if bindSupport, ok := di.WorkersContainerBind(bind); ok {
		bindSupport.HookUnregister()
	}

	c.binds.Delete(id)

	logger.Debug("Unregister bind",
		"type", bindItem.Type(),
		"id", bindItem.ID(),
	)

	if closer, ok := bind.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			logger.Debug("Unregister bind failed because close failed",
				"type", bindItem.Type(),
				"id", bindItem.ID(),
				"error", err.Error(),
			)
		}
	}

	c.itemStatusUpdate(bindItem, boggart.BindStatusRemoved)

	return nil
}

func (c *Component) unregisterBindsAll() error {
	for _, item := range c.BindItems() {
		if err := c.UnregisterBindByID(item.ID()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Component) Bind(id string) boggart.BindItem {
	if d, ok := c.binds.Load(id); ok {
		return d.(boggart.BindItem)
	}

	return nil
}

func (c *Component) BindItems() []boggart.BindItem {
	items := make([]boggart.BindItem, 0)

	c.binds.Range(func(_ interface{}, item interface{}) bool {
		items = append(items, item.(boggart.BindItem))
		return true
	})

	sort.Slice(items, func(i, j int) bool {
		if items[i].Type() == items[j].Type() {
			return items[i].ID() < items[j].ID()
		}

		return items[i].Type() < items[j].Type()
	})

	return items
}

func (c *Component) mqttOnConnectHandler(client mqtt.Component, restore bool) {
	for _, item := range c.BindItems() {
		topic := mqtt.Topic(c.config.String(boggart.ConfigMQTTTopicBindStatus)).Format(item.ID())

		err := client.PublishWithoutCache(context.Background(), topic, 1, true, strings.ToLower(item.Status().String()))
		if err != nil {
			c.itemLogger(item.Bind()).Error("Restore publish to MQTT failed", "topic", topic, "error", err.Error())
		}
	}
}

func (c *Component) itemStatusUpdate(item *BindItem, status boggart.BindStatus) {
	if ok := item.updateStatus(status); ok {
		topic := mqtt.Topic(c.config.String(boggart.ConfigMQTTTopicBindStatus)).Format(item.ID())
		name := status.String()
		payload := strings.ToLower(name)

		metricBindStatus.With("bind", item.ID(), "status", name).Inc()

		ctx := context.Background()

		// при закрытии шлем синхронно, что бы блочить операцию Close компонента
		if atomic.LoadInt32(&c.closing) != 0 {
			if err := c.mqtt.PublishWithoutCache(ctx, topic, 1, true, payload); err != nil {
				c.itemLogger(item.Bind()).Error("Publish to MQTT failed", "topic", topic, "error", err.Error())
			}
		} else {
			// FIXME: всегда нужно слать синхронно, так как асинхронность не гарантирует порядок
			// и статус иногда проскакивает не тот. Правильнее организовать внутри свою очередь,
			// чтобы не блочить работу с привязкой.

			// что бы публикация зарегистрировалась в общем списке и отображалась на странице mqtt для привязки
			if bindSupport, ok := di.MQTTContainerBind(item.Bind()); ok && bindSupport != nil {
				bindSupport.PublishAsync(ctx, topic, payload)
			} else {
				c.mqtt.PublishAsync(ctx, topic, 1, true, payload)
			}
		}
	}
}

func (c *Component) itemLogger(bind boggart.Bind) logging.Logger {
	if logger, ok := di.LoggerContainerBind(bind); ok {
		return logger
	}

	return c.logger
}

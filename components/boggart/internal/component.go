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
	_ "github.com/kihamo/boggart/components/boggart/bind/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
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
	workers     workers.Component

	routes  []dashboard.Route
	binds   sync.Map
	closing int32
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
	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	<-a.ReadyComponent(mqtt.ComponentName)
	c.mqtt = a.GetComponent(mqtt.ComponentName).(mqtt.Component)
	c.mqtt.OnConnectHandlerAdd(c.mqttOnConnectHandler)

	<-a.ReadyComponent(workers.ComponentName)
	c.workers = a.GetComponent(workers.ComponentName).(workers.Component)

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

	cfg, err := boggart.ValidateBindConfig(kind, map[string]interface{}{
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

		cfg, err := boggart.ValidateBindConfig(kind, d.Config)
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
	} else {
		if existsBind := c.Bind(id); existsBind != nil {
			return nil, errors.New("bind item with id " + id + " already exist")
		}
	}

	bindType, err := boggart.GetBindType(t)
	if err != nil {
		return nil, err
	}

	// register widget
	if widget, ok := bindType.(boggart.BindTypeHasWidgetAssetFS); ok {
		if fs := widget.WidgetAssetFS(); fs != nil {
			name := boggart.ComponentName + "-bind-" + t

			// templates
			render := c.dashboard.Renderer()
			if !render.IsRegisterNamespace(name) {
				if err := render.RegisterNamespace(name, fs); err != nil {
					return nil, err
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

	// mqtt container
	if bindSupport, ok := bind.(di.MQTTContainerSupport); ok {
		bindSupport.SetMQTT(di.NewMQTTContainer(bindItem, c.mqtt))
	}

	if runner, ok := bind.(boggart.BindRunner); ok {
		if err := runner.Run(); err != nil {
			return nil, err
		}
	}

	// mqtt subscribers
	if bindSupport, ok := di.MQTTContainerBind(bind); ok {
		// TODO: обвешать подписки враппером, что бы только в online можно было посылать
		for _, subscriber := range bindSupport.Subscribers() {
			if err := c.mqtt.SubscribeSubscriber(subscriber.Subscriber()); err != nil {
				c.itemStatusUpdate(bindItem, boggart.BindStatusUninitialized)
				return nil, err
			}
		}
	}

	tasks := make([]w.Task, 0)

	// probes container
	if bindSupport, ok := bind.(di.ProbesContainerSupport); ok {
		bindSupport.SetProbes(di.NewProbesContainer(
			bindItem,
			func(status boggart.BindStatus) {
				c.itemStatusUpdate(bindItem, status)
			}, func() error {
				_, err := c.RegisterBind(bindItem.ID(), bindItem.Bind(), bindItem.Type(), bindItem.Description(), bindItem.Tags(), bindItem.Config())
				return err
			}, func() error {
				return c.UnregisterBindByID(bindItem.ID())
			}))

		if probe := bindSupport.Probes().Readiness(); probe != nil {
			tasks = append(tasks, probe)
		} else {
			c.itemStatusUpdate(bindItem, boggart.BindStatusOnline)
		}

		if probe := bindSupport.Probes().Liveness(); probe != nil {
			tasks = append(tasks, probe)
		}
	} else {
		c.itemStatusUpdate(bindItem, boggart.BindStatusOnline)
	}

	// workers container
	if bindSupport, ok := bindItem.Bind().(di.WorkersContainerSupport); ok {
		bindSupport.SetWorkers(di.NewWorkersContainer(bindItem, c.workers))

		tasks = append(tasks, bindSupport.Workers().Tasks()...)
	}

	// register tasks
	for _, tsk := range tasks {
		c.workers.AddTask(tsk)
	}

	c.binds.Store(id, bindItem)

	c.logger.Debug("Register bind",
		"type", bindItem.Type(),
		"id", bindItem.ID(),
	)

	return bindItem, nil
}

func (c *Component) UnregisterBindByID(id string) error {
	d, ok := c.binds.Load(id)
	if !ok {
		return nil
	}

	bindItem := d.(*BindItem)

	c.itemStatusUpdate(bindItem, boggart.BindStatusRemoving)

	// unregister mqtt
	if bindSupport, ok := di.MQTTContainerBind(bindItem.Bind()); ok {
		// не блокирующее отписываемся, так как mqtt может быть не доступен
		go func() {
			for _, subscriber := range bindSupport.Subscribers() {
				if err := c.mqtt.UnsubscribeSubscriber(subscriber.Subscriber()); err != nil {
					c.logger.Error("Unregister bind failed because unsubscribe MQTT failed",
						"type", bindItem.Type(),
						"id", bindItem.ID(),
						"error", err.Error(),
					)
				}
			}

			if st := bindItem.Status(); st == boggart.BindStatusRemoving || st == boggart.BindStatusRemoved {
				bindSupport.SetClient(nil)
			}
		}()
	}

	// remove probes
	if bindSupport, ok := di.ProbesContainerBind(bindItem.Bind()); ok {
		if probe := bindSupport.Readiness(); probe != nil {
			c.workers.RemoveTask(probe)
		}

		if probe := bindSupport.Liveness(); probe != nil {
			c.workers.RemoveTask(probe)
		}
	}

	// workers
	if bindSupport, ok := di.WorkersContainerBind(bindItem.Bind()); ok {
		// remove tasks
		for _, tsk := range bindSupport.Tasks() {
			c.workers.RemoveTask(tsk)
		}
	}

	c.binds.Delete(id)

	c.logger.Debug("Unregister bind",
		"type", bindItem.Type(),
		"id", bindItem.ID(),
	)

	if closer, ok := bindItem.Bind().(io.Closer); ok {
		if err := closer.Close(); err != nil {
			c.logger.Debug("Unregister bind failed because close failed",
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
			c.logger.Error("Restore publish to MQTT failed", "topic", topic, "error", err.Error())
		}
	}
}

func (c *Component) itemStatusUpdate(item *BindItem, status boggart.BindStatus) {
	if ok := item.updateStatus(status); ok {
		topic := mqtt.Topic(c.config.String(boggart.ConfigMQTTTopicBindStatus)).Format(item.ID())
		payload := strings.ToLower(status.String())

		ctx := context.Background()

		// при закрытии шлем синхронно, что бы блочить операцию Close компонента
		if atomic.LoadInt32(&c.closing) != 0 {
			if err := c.mqtt.PublishWithoutCache(ctx, topic, 1, true, payload); err != nil {
				c.logger.Error("Publish to MQTT failed", "topic", topic, "error", err.Error())
			}
		} else {
			// что бы публикация зарегистрировалась в общем списке и отображалась на странице mqtt для привязки
			if bindSupport, ok := di.MQTTContainerBind(item.Bind()); ok && bindSupport != nil {
				bindSupport.PublishAsync(ctx, topic, payload)
			} else {
				c.mqtt.PublishAsync(ctx, topic, 1, true, payload)
			}
		}
	}
}

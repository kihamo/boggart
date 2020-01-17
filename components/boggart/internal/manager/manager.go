package manager

import (
	"context"
	"errors"
	"io"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

const (
	managerStatusNotReady = int64(iota)
	managerStatusReady
	managerStatusClose
)

type bindTask interface {
	w.Task
	SetName(string)
}

type Manager struct {
	status          int64
	storage         sync.Map
	dashboard       dashboard.Component
	i18n            i18n.Component
	mqtt            mqtt.Component
	workers         workers.Component
	logger          logging.Logger
	listeners       *manager.ListenersManager
	topicBindStatus mqtt.Topic
}

func NewManager(dashboard dashboard.Component,
	i18n i18n.Component,
	mqtt mqtt.Component,
	workers workers.Component,
	logger logging.Logger,
	listeners *manager.ListenersManager,
	topicBindStatus mqtt.Topic) *Manager {
	return &Manager{
		status:          managerStatusNotReady,
		dashboard:       dashboard,
		i18n:            i18n,
		mqtt:            mqtt,
		workers:         workers,
		logger:          logger,
		listeners:       listeners,
		topicBindStatus: topicBindStatus,
	}
}

func (m *Manager) Register(id string, bind boggart.Bind, t string, description string, tags []string, config interface{}) (boggart.BindItem, error) {
	if id == "" {
		id = uuid.New()
	} else {
		if existsBind := m.Bind(id); existsBind != nil {
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
			render := m.dashboard.Renderer()
			if !render.IsRegisterNamespace(name) {
				if err := render.RegisterNamespace(name, fs); err != nil {
					return nil, err
				}
			}

			// asset fs
			m.dashboard.RegisterAssetFS(name, fs)

			// i18n
			if m.i18n != nil {
				fs.Prefix = "locales"
				m.i18n.LoadLocaleFromFiles(name, i18n.FromAssetFS(fs))
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

	m.itemStatusUpdate(bindItem, boggart.BindStatusInitializing)

	// meta container
	if bindSupport, ok := bind.(di.MetaContainerSupport); ok {
		bindSupport.SetMeta(di.NewMetaContainer(bindItem))
	}

	// logger container
	if bindSupport, ok := bind.(di.LoggerContainerSupport); ok {
		bindSupport.SetLogger(di.NewLoggerContainer(logging.NewLazyLogger(m.logger, m.logger.Name()+"."+id)))
	}

	// mqtt container
	if bindSupport, ok := bind.(di.MQTTContainerSupport); ok {
		bindSupport.SetMQTT(di.NewMQTTContainer(bindItem, m.mqtt))
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
			if err := m.mqtt.SubscribeSubscriber(subscriber); err != nil {
				m.itemStatusUpdate(bindItem, boggart.BindStatusUninitialized)
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
				m.itemStatusUpdate(bindItem, status)
			}, func() error {
				_, err := m.Register(bindItem.ID(), bindItem.Bind(), bindItem.Type(), bindItem.Description(), bindItem.Tags(), bindItem.Config())
				return err
			}, func() error {
				return m.Unregister(bindItem.ID())
			}))

		if probe := bindSupport.Probes().Readiness(); probe != nil {
			tasks = append(tasks, probe)
		} else {
			m.itemStatusUpdate(bindItem, boggart.BindStatusOnline)
		}

		if probe := bindSupport.Probes().Liveness(); probe != nil {
			tasks = append(tasks, probe)
		}
	} else {
		m.itemStatusUpdate(bindItem, boggart.BindStatusOnline)
	}

	// workers container
	if bindSupport, ok := bindItem.Bind().(di.WorkersContainerSupport); ok {
		bindSupport.SetWorkers(di.NewWorkersContainer(bindItem))

		tasks = append(tasks, bindSupport.Workers().Tasks()...)

		// register listeners
		for _, listener := range bindSupport.Workers().Listeners() {
			m.listeners.AddListener(listener)
		}
	}

	// register tasks
	for _, tsk := range tasks {
		if tsk, ok := tsk.(bindTask); ok {
			tsk.SetName("bind-" + id + "-" + t + "-" + tsk.Name())
		}

		m.workers.AddTask(tsk)
	}

	m.storage.Store(id, bindItem)

	m.logger.Debug("Register bind",
		"type", bindItem.Type(),
		"id", bindItem.ID(),
	)

	return bindItem, nil
}

func (m *Manager) Unregister(id string) error {
	d, ok := m.storage.Load(id)
	if !ok {
		return nil
	}

	bindItem := d.(*BindItem)

	m.itemStatusUpdate(bindItem, boggart.BindStatusRemoving)

	// unregister mqtt
	if bindSupport, ok := di.MQTTContainerBind(bindItem.Bind()); ok {
		if err := m.mqtt.UnsubscribeSubscribers(bindSupport.Subscribers()); err != nil {
			m.logger.Error("Unregister bind failed because unsubscribe MQTT failed",
				"type", bindItem.Type(),
				"id", bindItem.ID(),
				"error", err.Error(),
			)
		}

		bindSupport.SetClient(nil)
	}

	// remove probes
	if bindSupport, ok := di.ProbesContainerBind(bindItem.Bind()); ok {
		if probe := bindSupport.Readiness(); probe != nil {
			m.workers.RemoveTask(probe)
		}

		if probe := bindSupport.Liveness(); probe != nil {
			m.workers.RemoveTask(probe)
		}
	}

	// workers
	if bindSupport, ok := di.WorkersContainerBind(bindItem.Bind()); ok {
		// remove tasks
		for _, tsk := range bindSupport.Tasks() {
			m.workers.RemoveTask(tsk)
		}

		// remove listeners
		for _, listener := range bindSupport.Listeners() {
			m.listeners.RemoveListener(listener)
		}
	}

	m.storage.Delete(id)

	m.logger.Debug("Unregister bind",
		"type", bindItem.Type(),
		"id", bindItem.ID(),
	)

	if closer, ok := bindItem.Bind().(io.Closer); ok {
		if err := closer.Close(); err != nil {
			m.logger.Debug("Unregister bind failed because close failed",
				"type", bindItem.Type(),
				"id", bindItem.ID(),
				"error", err.Error(),
			)
		}
	}

	m.itemStatusUpdate(bindItem, boggart.BindStatusRemoved)
	return nil
}

func (m *Manager) UnregisterAll() error {
	for _, item := range m.BindItems() {
		if err := m.Unregister(item.ID()); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Bind(id string) boggart.BindItem {
	if d, ok := m.storage.Load(id); ok {
		return d.(boggart.BindItem)
	}

	return nil
}

func (m *Manager) BindItems() BindItemsList {
	items := make([]boggart.BindItem, 0)

	m.storage.Range(func(_ interface{}, item interface{}) bool {
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

func (m *Manager) Describe(ch chan<- *snitch.Description) {
	m.storage.Range(func(_ interface{}, item interface{}) bool {
		if collector, ok := item.(*BindItem).Bind().(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (m *Manager) Collect(ch chan<- snitch.Metric) {
	m.storage.Range(func(_ interface{}, item interface{}) bool {
		if collector, ok := item.(*BindItem).Bind().(snitch.Collector); ok {
			collector.Collect(ch)
		}

		return true
	})
}

func (m *Manager) Ready() {
	atomic.SwapInt64(&m.status, managerStatusReady)
}

func (m *Manager) Close() error {
	old := atomic.SwapInt64(&m.status, managerStatusClose)
	if old != managerStatusClose {
		return m.UnregisterAll()
	}

	return nil
}

func (m *Manager) MQTTOnConnectHandler(client mqtt.Component, restore bool) {
	for _, item := range m.BindItems() {
		topic := m.topicBindStatus.Format(item.ID())
		err := client.PublishWithoutCache(context.Background(), topic, 1, true, strings.ToLower(item.Status().String()))

		if err != nil {
			m.logger.Error("Restore publish to MQTT failed", "topic", topic, "error", err.Error())
		}
	}
}

func (m *Manager) mqttPublish(topic mqtt.Topic, payload interface{}) {
	// при закрытии шлем синхронно, что бы блочить операцию Close компонента
	if atomic.LoadInt64(&m.status) == managerStatusClose {
		if err := m.mqtt.PublishWithoutCache(context.Background(), topic, 1, true, payload); err != nil {
			m.logger.Error("Publish to MQTT failed", "topic", topic, "error", err.Error())
		}
	} else {
		m.mqtt.PublishAsync(context.Background(), topic, 1, true, payload)
	}
}

func (m *Manager) itemStatusUpdate(item *BindItem, status boggart.BindStatus) {
	if ok := item.updateStatus(status); ok {
		m.mqttPublish(m.topicBindStatus.Format(item.id), strings.ToLower(status.String()))
	}
}

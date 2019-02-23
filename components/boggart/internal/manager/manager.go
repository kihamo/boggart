package manager

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

const (
	managerNotReady = int64(iota)
	managerReady

	MQTTPublishTopicBindStatus mqtt.Topic = boggart.ComponentName + "/bind/+/status"
)

type Manager struct {
	mutex sync.RWMutex

	ready     int64
	storage   *sync.Map
	dashboard dashboard.Component
	i18n      i18n.Component
	mqtt      mqtt.Component
	workers   workers.Component
	logger    logging.Logger
	listeners *manager.ListenersManager
}

func NewManager(dashboard dashboard.Component, i18n i18n.Component, mqtt mqtt.Component, workers workers.Component, logger logging.Logger, listeners *manager.ListenersManager) *Manager {
	return &Manager{
		ready:     managerNotReady,
		storage:   new(sync.Map),
		dashboard: dashboard,
		i18n:      i18n,
		mqtt:      mqtt,
		workers:   workers,
		logger:    logger,
		listeners: listeners,
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

	statusUpdate := m.itemStatusUpdate(bindItem)
	statusUpdate(boggart.BindStatusInitializing)

	// init logger
	if bindLogger, ok := bind.(boggart.BindLogger); ok {
		bindLogger.SetLogger(logging.NewLazyLogger(m.logger, m.logger.Name()+"."+id))
	}

	bind.SetStatusManager(bindItem.Status, m.bindStatusUpdate(bindItem))

	// register mqtt
	if mqttClient, ok := bind.(boggart.BindHasMQTTClient); ok {
		mqttClient.SetMQTTClient(m.mqtt)
	}

	if err := bind.Run(); err != nil {
		return nil, err
	}

	// TODO: обвешать подписки враппером, что бы только в online можно было посылать
	for _, subscriber := range bindItem.MQTTSubscribers() {
		if err := m.mqtt.SubscribeSubscriber(subscriber); err != nil {
			statusUpdate(boggart.BindStatusUninitialized)
			return nil, err
		}
	}

	// register tasks
	for _, task := range bindItem.Tasks() {
		m.workers.AddTask(task)
	}

	// register listeners
	for _, listener := range bindItem.Listeners() {
		m.listeners.AddListener(listener)
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

	statusUpdate := m.itemStatusUpdate(bindItem)
	statusUpdate(boggart.BindStatusRemoving)

	// unregister mqtt
	if err := m.mqtt.UnsubscribeSubscribers(bindItem.MQTTSubscribers()); err != nil {
		statusUpdate(boggart.BindStatusUnknown)
		return err
	}

	if mqttClient, ok := bindItem.Bind().(boggart.BindHasMQTTClient); ok {
		mqttClient.SetMQTTClient(nil)
	}

	// remove tasks
	for _, task := range bindItem.Tasks() {
		m.workers.RemoveTask(task)
	}

	// remove listeners
	for _, listener := range bindItem.Listeners() {
		m.listeners.RemoveListener(listener)
	}

	m.storage.Delete(id)

	m.logger.Debug("Unregister bind",
		"type", bindItem.Type(),
		"id", bindItem.ID(),
	)

	if closer, ok := bindItem.Bind().(boggart.BindCloser); ok {
		if err := closer.Close(); err != nil {
			statusUpdate(boggart.BindStatusUnknown)
			return err
		}
	}

	statusUpdate(boggart.BindStatusRemoved)
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
	if !m.IsReady() {
		atomic.StoreInt64(&m.ready, managerReady)
		m.listeners.AsyncTrigger(context.TODO(), boggart.BindEventManagerReady)
	}
}

func (m *Manager) IsReady() bool {
	return atomic.LoadInt64(&m.ready) == managerReady
}

func (m *Manager) itemStatusUpdate(item *BindItem) boggart.BindStatusSetter {
	return func(status boggart.BindStatus) {
		if status == item.Status() {
			return
		}

		item.updateStatus(status)

		m.mqtt.PublishAsync(
			context.Background(),
			MQTTPublishTopicBindStatus.Format(mqtt.NameReplace(item.id)),
			1,
			true,
			strings.ToLower(status.String()))
	}
}

func (m *Manager) bindStatusUpdate(item *BindItem) boggart.BindStatusSetter {
	return func(status boggart.BindStatus) {
		switch status {
		// allow statuses
		case boggart.BindStatusOnline, boggart.BindStatusOffline, boggart.BindStatusUnknown, boggart.BindStatusRemoved:
			if status == item.Status() {
				return
			}

			item.updateStatus(status)

			m.mqtt.PublishAsync(
				context.Background(),
				MQTTPublishTopicBindStatus.Format(mqtt.NameReplace(item.id)),
				1,
				true,
				strings.ToLower(status.String()))

		default:
			m.logger.Error("Deny change status", "status", status.String(), "item-id", item.id)
		}
	}
}

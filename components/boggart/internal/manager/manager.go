package manager

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

const (
	managerNotReady = int64(iota)
	managerReady
)

type Manager struct {
	mutex sync.RWMutex

	ready     int64
	storage   *sync.Map
	mqtt      mqtt.Component
	workers   workers.Component
	listeners *manager.ListenersManager
}

func NewManager(mqtt mqtt.Component, workers workers.Component, listeners *manager.ListenersManager) *Manager {
	return &Manager{
		ready:     managerNotReady,
		storage:   new(sync.Map),
		mqtt:      mqtt,
		workers:   workers,
		listeners: listeners,
	}
}

func (m *Manager) Register(bind boggart.Bind, t string, description string, tags []string, config interface{}) (string, error) {
	id := uuid.New()
	err := m.RegisterWithID(id, bind, t, description, tags, config)
	return id, err
}

func (m *Manager) RegisterWithID(id string, bind boggart.Bind, t string, description string, tags []string, config interface{}) error {
	if id == "" {
		id = uuid.New()
	}

	bindItem := &BindItem{
		bind:        bind,
		id:          id,
		t:           t,
		description: description,
		tags:        tags,
		config:      config,
	}
	m.storage.Store(id, bindItem)
	m.listeners.AsyncTrigger(context.TODO(), boggart.BindEventDeviceRegister, bind, id)

	// register mqtt
	if mqttClient, ok := bind.(boggart.BindHasMQTTClient); ok {
		mqttClient.SetMQTTClient(m.mqtt)
	}

	for _, subscriber := range bindItem.MQTTSubscribers() {
		if err := m.mqtt.SubscribeSubscriber(subscriber); err != nil {
			return err
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

	return nil
}

func (m *Manager) Unregister(id string) error {
	d, ok := m.storage.Load(id)
	if !ok {
		return nil
	}

	bindItem := d.(*BindItem)

	// unregister mqtt
	if err := m.mqtt.UnsubscribeSubscribers(bindItem.MQTTSubscribers()); err != nil {
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

	if closer, ok := bindItem.Bind().(boggart.BindCloser); ok {
		return closer.Close()
	}

	return nil
}

func (m *Manager) Bind(id string) boggart.Device {
	if d, ok := m.storage.Load(id); ok {
		return d.(boggart.Device)
	}

	return nil
}

func (m *Manager) BindItems() BindItemsList {
	items := make([]*BindItem, 0)

	m.storage.Range(func(key interface{}, item interface{}) bool {
		items = append(items, item.(*BindItem))
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
		m.listeners.AsyncTrigger(context.TODO(), boggart.BindEventDevicesManagerReady)
	}
}

func (m *Manager) IsReady() bool {
	return atomic.LoadInt64(&m.ready) == managerReady
}

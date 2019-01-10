package internal

import (
	"context"
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
	devicesManagerNotReady = int64(iota)
	devicesManagerReady
)

type DevicesManager struct {
	mutex sync.RWMutex

	ready     int64
	storage   *sync.Map
	mqtt      mqtt.Component
	workers   workers.Component
	listeners *manager.ListenersManager
}

func NewDevicesManager(mqtt mqtt.Component, workers workers.Component, listeners *manager.ListenersManager) *DevicesManager {
	return &DevicesManager{
		ready:     devicesManagerNotReady,
		storage:   new(sync.Map),
		mqtt:      mqtt,
		workers:   workers,
		listeners: listeners,
	}
}

func (m *DevicesManager) Register(device boggart.DeviceBind, t string, description string, tags []string, config map[string]interface{}) string {
	id := uuid.New()
	m.RegisterWithID(id, device, t, description, tags, config)
	return id
}

func (m *DevicesManager) RegisterWithID(id string, bind boggart.DeviceBind, t string, description string, tags []string, config map[string]interface{}) {
	if id == "" {
		id = uuid.New()
	}

	m.storage.Store(id, &Device{
		bind:        bind,
		id:          id,
		t:           t,
		description: description,
		tags:        tags,
		config:      config,
	})
	m.listeners.AsyncTrigger(context.TODO(), boggart.DeviceEventDeviceRegister, bind, id)

	if mqttClient, ok := bind.(boggart.DeviceBindHasMQTTClient); ok {
		mqttClient.SetMQTTClient(m.mqtt)
	}

	if subs, ok := bind.(boggart.DeviceBindHasMQTTSubscribers); ok {
		m.mqtt.SubscribeSubscribers(subs.MQTTSubscribers())
	}

	if tasks, ok := bind.(boggart.DeviceBindHasTasks); ok {
		for _, task := range tasks.Tasks() {
			m.workers.AddTask(task)
		}
	}

	if listeners, ok := bind.(boggart.DeviceBindHasListeners); ok {
		for _, listener := range listeners.Listeners() {
			m.listeners.AddListener(listener)
		}
	}
}

func (m *DevicesManager) Device(id string) boggart.Device {
	if d, ok := m.storage.Load(id); ok {
		return d.(boggart.Device)
	}

	return nil
}

func (m *DevicesManager) Devices() map[string]boggart.Device {
	devices := make(map[string]boggart.Device, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		devices[key.(string)] = device.(boggart.Device)
		return true
	})

	return devices
}

func (m *DevicesManager) Describe(ch chan<- *snitch.Description) {
	m.storage.Range(func(_ interface{}, device interface{}) bool {
		if collector, ok := device.(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (m *DevicesManager) Collect(ch chan<- snitch.Metric) {
	m.storage.Range(func(_ interface{}, device interface{}) bool {
		if collector, ok := device.(snitch.Collector); ok {
			collector.Collect(ch)
		}

		return true
	})
}

func (m *DevicesManager) Ready() {
	if !m.IsReady() {
		atomic.StoreInt64(&m.ready, devicesManagerReady)
		m.listeners.AsyncTrigger(context.TODO(), boggart.DeviceEventDevicesManagerReady)
	}
}

func (m *DevicesManager) IsReady() bool {
	return atomic.LoadInt64(&m.ready) == devicesManagerReady
}

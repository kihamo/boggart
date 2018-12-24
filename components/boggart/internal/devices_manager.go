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

func (m *DevicesManager) Register(device boggart.Device) string {
	id := uuid.New()
	m.RegisterWithID(id, device)
	return id
}

func (m *DevicesManager) RegisterWithID(id string, device boggart.Device) {
	m.storage.Store(id, device)
	m.listeners.AsyncTrigger(context.TODO(), boggart.DeviceEventDeviceRegister, device, id)

	if mqttClient, ok := device.(boggart.DeviceHasMQTTClient); ok {
		mqttClient.SetMQTTClient(m.mqtt)
	}

	if subs, ok := device.(boggart.DeviceHasMQTTSubscribers); ok {
		m.mqtt.SubscribeSubscribers(subs.MQTTSubscribers())
	}

	if tasks, ok := device.(boggart.DeviceHasTasks); ok {
		for _, task := range tasks.Tasks() {
			m.workers.AddTask(task)
		}
	}

	if listeners, ok := device.(boggart.DeviceHasListeners); ok {
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

func (m *DevicesManager) DevicesByTypes(types []boggart.DeviceType) map[string]boggart.Device {
	if len(types) == 0 {
		return m.Devices()
	}

	devices := make(map[string]boggart.Device, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		d := device.(boggart.Device)

		for _, t1 := range d.Types() {
			for _, t2 := range types {
				if t1.String() == t2.String() {
					devices[key.(string)] = d
					return true
				}
			}
		}

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

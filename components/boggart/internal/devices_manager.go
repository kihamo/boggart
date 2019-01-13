package internal

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

func (m *DevicesManager) Register(device boggart.DeviceBind, t string, description string, tags []string, config interface{}) (string, error) {
	id := uuid.New()
	err := m.RegisterWithID(id, device, t, description, tags, config)
	return id, err
}

func (m *DevicesManager) RegisterWithID(id string, bind boggart.DeviceBind, t string, description string, tags []string, config interface{}) error {
	if id == "" {
		id = uuid.New()
	}

	device := &DeviceItem{
		bind:        bind,
		id:          id,
		t:           t,
		description: description,
		tags:        tags,
		config:      config,
	}
	m.storage.Store(id, device)
	m.listeners.AsyncTrigger(context.TODO(), boggart.DeviceEventDeviceRegister, bind, id)

	// register mqtt
	if mqttClient, ok := bind.(boggart.DeviceBindHasMQTTClient); ok {
		mqttClient.SetMQTTClient(m.mqtt)
	}

	for _, subscriber := range device.MQTTSubscribers() {
		if err := m.mqtt.SubscribeSubscriber(subscriber); err != nil {
			return err
		}
	}

	// register tasks
	for _, task := range device.Tasks() {
		m.workers.AddTask(task)
	}

	// register listeners
	for _, listener := range device.Listeners() {
		m.listeners.AddListener(listener)
	}

	return nil
}

func (m *DevicesManager) Unregister(id string) error {
	d, ok := m.storage.Load(id)
	if !ok {
		return nil
	}

	device := d.(*DeviceItem)

	// unregister mqtt
	if err := m.mqtt.UnsubscribeSubscribers(device.MQTTSubscribers()); err != nil {
		return err
	}

	if mqttClient, ok := device.Bind().(boggart.DeviceBindHasMQTTClient); ok {
		mqttClient.SetMQTTClient(nil)
	}

	// remove tasks
	for _, task := range device.Tasks() {
		m.workers.RemoveTask(task)
	}

	// remove listeners
	for _, listener := range device.Listeners() {
		m.listeners.RemoveListener(listener)
	}

	m.storage.Delete(id)

	if closer, ok := device.Bind().(boggart.DeviceBindCloser); ok {
		return closer.Close()
	}

	return nil
}

func (m *DevicesManager) Device(id string) boggart.Device {
	if d, ok := m.storage.Load(id); ok {
		return d.(boggart.Device)
	}

	return nil
}

func (m *DevicesManager) Devices() []boggart.Device {
	devices := make([]boggart.Device, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		devices = append(devices, device.(boggart.Device))
		return true
	})

	sort.Slice(devices, func(i, j int) bool {
		if devices[i].Type() == devices[j].Type() {
			return devices[i].ID() < devices[j].ID()
		}

		return devices[i].Type() < devices[j].Type()
	})

	return devices
}

func (m *DevicesManager) Describe(ch chan<- *snitch.Description) {
	m.storage.Range(func(_ interface{}, device interface{}) bool {
		if collector, ok := device.(*DeviceItem).Bind().(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (m *DevicesManager) Collect(ch chan<- snitch.Metric) {
	m.storage.Range(func(_ interface{}, device interface{}) bool {
		if collector, ok := device.(*DeviceItem).Bind().(snitch.Collector); ok {
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

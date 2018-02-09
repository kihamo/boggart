package internal

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

type DeviceManager struct {
	storage *sync.Map
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		new(sync.Map),
	}
}

func (m *DeviceManager) Register(id string, device boggart.Device) {
	m.storage.Store(id, device)
}

func (m *DeviceManager) Device(id string) boggart.Device {
	if d, ok := m.storage.Load(id); ok {
		return d.(boggart.Device)
	}

	return nil
}

func (m *DeviceManager) Devices() map[string]boggart.Device {
	devices := make(map[string]boggart.Device, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		devices[key.(string)] = device.(boggart.Device)
		return true
	})

	return devices
}

func (m *DeviceManager) Describe(ch chan<- *snitch.Description) {
	m.storage.Range(func(key interface{}, device interface{}) bool {
		if collector, ok := device.(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (m *DeviceManager) Collect(ch chan<- snitch.Metric) {
	m.storage.Range(func(key interface{}, device interface{}) bool {
		if collector, ok := device.(snitch.Collector); ok {
			collector.Collect(ch)
		}

		return true
	})
}

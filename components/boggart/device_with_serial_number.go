package boggart

import (
	"sync"
)

type DeviceWithSerialNumber struct {
	DeviceBase

	mutex        sync.RWMutex
	serialNumber string
}

func (d *DeviceWithSerialNumber) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *DeviceWithSerialNumber) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

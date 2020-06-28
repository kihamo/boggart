package z_stack

import (
	"encoding/hex"
	"sync/atomic"

	a "github.com/kihamo/boggart/atomic"
)

type Device struct {
	networkAddress uint32
	ieeeAddress    *a.Bytes
	capabilities   uint32
	deviceType     uint32
}

func NewDevice() *Device {
	return &Device{
		ieeeAddress: a.NewBytes(),
	}
}

func (d *Device) NetworkAddress() uint16 {
	return uint16(atomic.LoadUint32(&d.networkAddress))
}

func (d *Device) SetNetworkAddress(address uint16) {
	atomic.StoreUint32(&d.networkAddress, uint32(address))
}

func (d *Device) IEEEAddress() []byte {
	return d.ieeeAddress.Load()
}

func (d *Device) SetIEEEAddress(address []byte) {
	d.ieeeAddress.Set(address)
}

func (d *Device) IEEEAddressAsString() string {
	return hex.EncodeToString(d.IEEEAddress())
}

func (d *Device) Capabilities() uint8 {
	return uint8(atomic.LoadUint32(&d.capabilities))
}

func (d *Device) SetCapabilities(capabilities uint8) {
	atomic.StoreUint32(&d.capabilities, uint32(capabilities))
}

func (d *Device) DeviceType() uint8 {
	return uint8(atomic.LoadUint32(&d.deviceType))
}

func (d *Device) DeviceTypeAsString() string {
	switch d.DeviceType() {
	case DeviceTypeCoordinator:
		return "coordinator"
	case DeviceTypeRouter:
		return "router"
	case DeviceTypeEndDevice:
		return "end device"
	}

	return "none"
}

func (d *Device) SetDeviceType(deviceType uint8) {
	atomic.StoreUint32(&d.deviceType, uint32(deviceType))
}

func (c *Client) deviceAdd(device *Device) {
	c.devices.Store(device.NetworkAddress(), device)
}

func (c *Client) deviceRemove(networkAddress uint16) {
	c.devices.Delete(networkAddress)
}

func (c *Client) Devices() []*Device {
	devices := make([]*Device, 0)

	c.devices.Range(func(key, value interface{}) bool {
		devices = append(devices, value.(*Device))
		return true
	})

	return devices
}

func (c *Client) Device(networkAddress uint16) *Device {
	value, ok := c.devices.Load(networkAddress)
	if !ok {
		return nil
	}

	return value.(*Device)
}

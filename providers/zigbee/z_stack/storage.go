package zstack

import (
	"github.com/kihamo/boggart/providers/zigbee/z_stack/model"
)

func (c *Client) deviceAdd(device *model.Device) {
	c.devices.Store(device.NetworkAddress(), device)
}

func (c *Client) deviceRemove(networkAddress uint16) {
	c.devices.Delete(networkAddress)
}

func (c *Client) Devices() []*model.Device {
	devices := make([]*model.Device, 0)

	c.devices.Range(func(key, value interface{}) bool {
		devices = append(devices, value.(*model.Device))
		return true
	})

	return devices
}

func (c *Client) Device(networkAddress uint16) *model.Device {
	value, ok := c.devices.Load(networkAddress)
	if !ok {
		return nil
	}

	return value.(*model.Device)
}

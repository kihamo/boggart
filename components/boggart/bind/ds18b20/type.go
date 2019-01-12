package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	device := &Bind{
		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(config.Address)

	return device, nil
}

package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	device := &Bind{
		provider:  softvideo.NewClient(config.Login, config.Password),
		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(config.Login)

	return device, nil
}

package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	device := &Bind{
		provider:        softvideo.NewClient(config.Login, config.Password, config.Debug),
		lastValue:       -1,
		updaterInterval: config.UpdaterInterval,
	}
	device.SetSerialNumber(config.Login)

	return device, nil
}

package softvideo

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	updaterInterval, err := time.ParseDuration(config.UpdaterInterval)
	if err != nil {
		return nil, err
	}

	device := &Bind{
		provider:        softvideo.NewClient(config.Login, config.Password),
		lastValue:       -1,
		updaterInterval: updaterInterval,
	}
	device.Init()
	device.SetSerialNumber(config.Login)

	return device, nil
}

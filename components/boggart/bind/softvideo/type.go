package softvideo

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (bind boggart.DeviceBind, err error) {
	config := c.(*Config)

	var updaterInterval time.Duration

	if config.UpdaterInterval != "" {
		updaterInterval, err = time.ParseDuration(config.UpdaterInterval)
		if err != nil {
			return nil, err
		}
	} else {
		updaterInterval = DefaultUpdaterInterval
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

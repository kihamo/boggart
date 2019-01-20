package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	device := &Bind{
		lastValue:        -1,
		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
		updaterInterval:  config.UpdaterInterval,
	}

	device.SetSerialNumber(config.Address)

	return device, nil
}

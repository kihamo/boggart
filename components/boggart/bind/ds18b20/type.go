package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		temperature:      atomic.NewFloat32Null(),
		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
		updaterInterval:  config.UpdaterInterval,
	}

	bind.SetSerialNumber(config.Address)

	return bind, nil
}

package samsung_tizen

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	livenessInterval, err := time.ParseDuration(config.LivenessInterval)
	if err != nil {
		return nil, err
	}

	livenessTimeout, err := time.ParseDuration(config.LivenessTimeout)
	if err != nil {
		return nil, err
	}

	device := &Bind{
		client:           tv.NewApiV2(config.Host),
		livenessInterval: livenessInterval,
		livenessTimeout:  livenessTimeout,
	}
	device.Init()

	return device, nil
}

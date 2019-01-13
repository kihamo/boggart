package lg_webos

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
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
		livenessInterval: livenessInterval,
		livenessTimeout:  livenessTimeout,
		host:             config.Host,
		key:              config.Key,
	}
	device.Init()

	return device, nil
}

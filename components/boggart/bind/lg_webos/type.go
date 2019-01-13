package lg_webos

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	device := &Bind{
		host:             config.Host,
		key:              config.Key,
		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}
	device.Init()

	return device, nil
}

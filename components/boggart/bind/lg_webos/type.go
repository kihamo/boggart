package lg_webos

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		host:             config.Host,
		key:              config.Key,
		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}

	return bind, nil
}

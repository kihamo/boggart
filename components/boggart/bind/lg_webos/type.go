package lg_webos

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type Type struct {
	dashboard.Handler
}

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

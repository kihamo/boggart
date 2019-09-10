package samsung_tizen

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		client:           tv.NewApiV2(config.Host),
		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}

	return bind, nil
}

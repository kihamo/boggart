package sun

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
		config: config,
	}
	bind.SetSerialNumber(config.Name)

	return bind, nil
}

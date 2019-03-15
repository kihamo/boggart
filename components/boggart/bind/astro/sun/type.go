package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:   config,
		sunrise:  atomic.NewTimeNull(),
		sunset:   atomic.NewTimeNull(),
		dayLight: atomic.NewDuration(),
	}
	bind.SetSerialNumber(config.Name)

	return bind, nil
}

package premiergc

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/premiergc"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:   config,
		provider: premiergc.New(config.Login, config.Password, config.Debug),
	}

	return bind, nil
}

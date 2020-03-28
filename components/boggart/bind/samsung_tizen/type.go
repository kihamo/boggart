package tizen

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		client: tv.NewAPIv2(config.Host),
	}

	return bind, nil
}

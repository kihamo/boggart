package esphome

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:   config,
		provider: native_api.New(config.Address, config.Password).WithClientID("Boggart bind"),
	}

	return bind, nil
}

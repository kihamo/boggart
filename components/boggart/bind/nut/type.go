package nut

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		host:            config.Host,
		ups:             config.UPS,
		username:        config.Username,
		password:        config.Password,
		updaterInterval: config.UpdaterInterval,
		variables:       make(map[string]interface{}),
	}

	return bind, nil
}

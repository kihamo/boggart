package nut

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	device := &Bind{
		host:            config.Host,
		ups:             config.UPS,
		username:        config.Username,
		password:        config.Password,
		updaterInterval: config.UpdaterInterval,
		variables:       make(map[string]interface{}, 0),
	}
	device.Init()

	return device, nil
}

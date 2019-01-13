package nut

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	updaterInterval, err := time.ParseDuration(config.UpdaterInterval)
	if err != nil {
		return nil, err
	}

	device := &Bind{
		host:            config.Host,
		ups:             config.UPS,
		username:        config.Username,
		password:        config.Password,
		updaterInterval: updaterInterval,
		variables:       make(map[string]interface{}, 0),
	}
	device.Init()

	return device, nil
}

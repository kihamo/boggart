package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		device: vacuum.New(config.Host, config.Token),
	}

	if config.PacketsCounter > 0 {
		bind.device.Client().SetPacketsCounter(config.PacketsCounter)
	}

	return bind, nil
}

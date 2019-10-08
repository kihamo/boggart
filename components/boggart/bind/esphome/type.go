package esphome

import (
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		provider: native_api.New(config.Address, config.Password).
			WithClientID("Boggart bind").
			WithDebug(config.Debug),
		otaAddress: config.Address,
	}

	if config.OTAPort > 0 {
		bind.otaAddress += strconv.FormatUint(config.OTAPort, 10)
	}

	return bind, nil
}

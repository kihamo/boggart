package nativeapi

import (
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome"
	api "github.com/kihamo/boggart/providers/esphome/native_api"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	otaAddress := config.Address
	if config.OTAPort > 0 {
		otaAddress += strconv.FormatUint(config.OTAPort, 10)
	}

	bind := &Bind{
		config: config,
		provider: api.New(config.Address, config.Password).
			WithClientID("Boggart bind").
			WithDebug(config.Debug),
		ota: esphome.NewOTA(otaAddress, config.OTAPassword),
	}

	return bind, nil
}

package rm

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*ConfigRM)

	var provider broadlink.Device

	switch config.Model {
	case "rm3mini":
		provider = broadlink.NewRMMini(config.MAC.HardwareAddr, config.Host)

	case "rm2proplus":
		provider = broadlink.NewRM2ProPlus3(config.MAC.HardwareAddr, config.Host)

	default:
		return nil, errors.New("unknown model " + config.Model)
	}

	provider.SetTimeout(config.ConnectionTimeout)

	bind := &Bind{
		provider:        provider,
		mac:             config.MAC.HardwareAddr,
		host:            config.Host,
		captureDuration: config.CaptureDuration,

		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}
	bind.SetSerialNumber(config.MAC.String())

	return bind, nil
}

package rm

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/broadlink"
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

	sn := config.MAC.String()

	config.TopicCapture = config.TopicCapture.Format(sn)
	config.TopicCaptureState = config.TopicCaptureState.Format(sn)
	config.TopicIR = config.TopicIR.Format(sn)
	config.TopicIRCount = config.TopicIRCount.Format(sn)
	config.TopicIRCapture = config.TopicIRCapture.Format(sn)
	config.TopicRF315mhz = config.TopicRF315mhz.Format(sn)
	config.TopicRF315mhzCount = config.TopicRF315mhzCount.Format(sn)
	config.TopicRF315mhzCapture = config.TopicRF315mhzCapture.Format(sn)
	config.TopicRF433mhz = config.TopicRF433mhz.Format(sn)
	config.TopicRF433mhzCount = config.TopicRF433mhzCount.Format(sn)
	config.TopicRF433mhzCapture = config.TopicRF433mhzCapture.Format(sn)

	provider.SetTimeout(config.ConnectionTimeout)

	bind := &Bind{
		config:   config,
		provider: provider,
	}
	bind.SetSerialNumber(sn)

	return bind, nil
}

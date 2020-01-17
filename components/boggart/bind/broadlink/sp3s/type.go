package sp3s

import (
	"errors"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/broadlink"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	var provider *broadlink.SP3S

	switch config.Model {
	case "sp3seu":
		provider = broadlink.NewSP3SEU(config.MAC.HardwareAddr, config.Host)

	case "sp3sus":
		provider = broadlink.NewSP3SUS(config.MAC.HardwareAddr, config.Host)

	default:
		return nil, errors.New("unknown model " + config.Model)
	}

	sn := config.MAC.String()

	config.TopicState = config.TopicState.Format(sn)
	config.TopicPower = config.TopicPower.Format(sn)
	config.TopicSet = config.TopicSet.Format(sn)

	provider.SetTimeout(config.ConnectionTimeout)

	bind := &Bind{
		config:   config,
		provider: provider,
		state:    atomic.NewBoolNull(),
		power:    atomic.NewFloat32Null(),
	}

	return bind, nil
}

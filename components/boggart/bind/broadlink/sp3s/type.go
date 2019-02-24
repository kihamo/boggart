package sp3s

import (
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	localAddr, err := broadlink.LocalAddr()
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   config.IP.IP,
		Port: broadlink.DevicePort,
	}

	var provider *broadlink.SP3S

	switch config.Model {
	case "sp3seu":
		provider = broadlink.NewSP3SEU(config.MAC.HardwareAddr, ip, *localAddr)

	case "sp3sus":
		provider = broadlink.NewSP3SUS(config.MAC.HardwareAddr, ip, *localAddr)

	default:
		return nil, errors.New("unknown model " + config.Model)
	}

	provider.SetTimeout(config.ConnectionTimeout)

	bind := &Bind{
		provider:        provider,
		updaterInterval: config.UpdaterInterval,
		state:           atomic.NewBoolNull(),
		power:           atomic.NewFloat32Null(),
	}
	bind.SetSerialNumber(config.MAC.String())

	return bind, nil
}

package rm

import (
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*ConfigRM)

	localAddr, err := broadlink.LocalAddr()
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   config.IP.IP,
		Port: broadlink.DevicePort,
	}

	var provider interface{}

	switch config.Model {
	case "rm3mini":
		provider = broadlink.NewRMMini(config.MAC.HardwareAddr, ip, *localAddr)

	case "rm2proplus":
		provider = broadlink.NewRM2ProPlus3(config.MAC.HardwareAddr, ip, *localAddr)

	default:
		return nil, errors.New("unknown model " + config.Model)
	}

	bind := &Bind{
		provider:        provider,
		mac:             config.MAC.HardwareAddr,
		ip:              ip,
		captureDuration: config.CaptureDuration,

		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}
	bind.SetSerialNumber(config.MAC.String())

	return bind, nil
}

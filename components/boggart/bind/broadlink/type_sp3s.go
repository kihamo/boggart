package broadlink

import (
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type TypeSP3S struct{}

func (t TypeSP3S) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*ConfigSP3S)

	mac, err := net.ParseMAC(config.MAC)
	if err != nil {
		return nil, err
	}

	localAddr, err := broadlink.LocalAddr()
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   net.ParseIP(config.IP),
		Port: broadlink.DevicePort,
	}

	var provider *broadlink.SP3S

	switch config.Model {
	case "sp3seu":
		provider = broadlink.NewSP3SEU(mac, ip, *localAddr)

	case "sp3sus":
		provider = broadlink.NewSP3SUS(mac, ip, *localAddr)

	default:
		return nil, errors.New("unknown model " + config.Model)
	}

	device := &BindSP3S{
		provider:        provider,
		state:           0,
		power:           -1,
		updaterInterval: config.UpdaterInterval,
	}
	device.Init()
	device.SetSerialNumber(mac.String())

	return device, nil
}

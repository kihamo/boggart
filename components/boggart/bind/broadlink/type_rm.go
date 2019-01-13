package broadlink

import (
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type TypeRM struct{}

func (t TypeRM) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*ConfigRM)

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

	var provider interface{}

	switch config.Model {
	case "rm3mini":
		provider = broadlink.NewRMMini(mac, ip, *localAddr)

	case "rm2proplus":
		provider = broadlink.NewRM2ProPlus3(mac, ip, *localAddr)

	default:
		return nil, errors.New("unknown model " + config.Model)
	}

	device := &BindRM{
		provider:        provider,
		mac:             mac,
		ip:              ip,
		captureDuration: config.CaptureDuration,
	}
	device.Init()
	device.SetSerialNumber(mac.String())

	// TODO: check open UDP port
	device.UpdateStatus(boggart.DeviceStatusOnline)

	return device, nil
}

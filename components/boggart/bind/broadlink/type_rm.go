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

	device := &BindRM{
		provider:        provider,
		mac:             config.MAC.HardwareAddr,
		ip:              ip,
		captureDuration: config.CaptureDuration,
	}
	device.Init()
	device.SetSerialNumber(config.MAC.String())

	// TODO: check open (ping) UDP port
	device.UpdateStatus(boggart.DeviceStatusOnline)

	return device, nil
}

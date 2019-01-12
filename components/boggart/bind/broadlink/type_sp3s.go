package broadlink

import (
	"net"
	"time"

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

	updaterInterval, err := time.ParseDuration(config.UpdaterInterval)
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

	device := &BindSP3S{
		provider:        broadlink.NewSP3S(mac, ip, *localAddr),
		state:           0,
		power:           -1,
		updaterInterval: updaterInterval,
	}
	device.Init()
	device.SetSerialNumber(mac.String())

	return device, nil
}

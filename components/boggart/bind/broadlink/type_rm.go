package broadlink

import (
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

	mac, err := net.ParseMAC(config.MAC)
	if err != nil {
		return nil, err
	}

	ip := net.UDPAddr{
		IP:   net.ParseIP(config.IP),
		Port: broadlink.DevicePort,
	}

	device := &BindRM{
		provider: broadlink.NewRMProPlus(mac, ip, *localAddr),
	}
	device.Init()
	device.SetSerialNumber(mac.String())

	return device, nil
}

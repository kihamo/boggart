package wol

import (
	"errors"
	"net"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}

func (b *Bind) WOL(mac net.HardwareAddr, ip net.IP, subnet net.IP) error {
	if mac == nil {
		return errors.New("MAC isn't set")
	}

	var broadcastAddress net.IP

	if ip != nil && subnet != nil {
		broadcastAddress = net.IP{0, 0, 0, 0}
		for i := 0; i < 4; i++ {
			broadcastAddress[i] = (ip[i] & subnet[i]) | ^subnet[i]
		}
	} else {
		broadcastAddress = net.IP{255, 255, 255, 255}
	}

	return wol.MagicWake(mac.String(), broadcastAddress.String())
}

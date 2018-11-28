package broadlink

import (
	"net"

	"github.com/kihamo/boggart/components/boggart/providers/broadlink/internal"
)

type RMProPlus struct {
	*RM3Mini
}

func NewRMProPlus(mac net.HardwareAddr, addr, iface net.UDPAddr) *RMProPlus {
	return &RMProPlus{
		RM3Mini: &RM3Mini{
			Device: internal.NewDevice(KindRMProPlus, mac, addr, iface),
		},
	}
}

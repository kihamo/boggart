package wol

import (
	"net"

	w "github.com/ghthor/gowol"
)

func WakeUp(mac, ip, subnet string) error {
	netIP := net.ParseIP(ip)
	netSubnet := net.ParseIP(subnet)

	broadcastAddress := net.IP{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		broadcastAddress[i] = (netIP[i] & netSubnet[i]) | ^netSubnet[i]
	}

	return w.MagicWake(mac, broadcastAddress.String())
}

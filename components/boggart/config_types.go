package boggart

import (
	"net"
)

type IP struct {
	net.IP
}

func (ip IP) MarshalText() ([]byte, error) {
	return []byte(ip.String()), nil
}

type HardwareAddr struct {
	net.HardwareAddr
}

func (addr HardwareAddr) MarshalText() ([]byte, error) {
	return []byte(addr.String()), nil
}

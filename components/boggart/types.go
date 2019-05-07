package boggart

import (
	"net"
	"net/url"
)

type IP struct {
	net.IP
}

func (ip IP) MarshalText() ([]byte, error) {
	return []byte(ip.String()), nil
}

type IPMask struct {
	net.IPMask
}

func (mask *IPMask) UnmarshalText(text []byte) error {
	ip := &net.IP{}

	err := ip.UnmarshalText(text)
	if err == nil {
		*mask = IPMask{
			IPMask: net.IPMask(ip.To4()),
		}
	}

	return err
}

type HardwareAddr struct {
	net.HardwareAddr
}

func (addr HardwareAddr) MarshalText() ([]byte, error) {
	return []byte(addr.String()), nil
}

func (addr *HardwareAddr) UnmarshalText(text []byte) error {
	mac, err := net.ParseMAC(string(text))

	if err == nil {
		*addr = HardwareAddr{
			HardwareAddr: mac,
		}
	}

	return err
}

type URL struct {
	url.URL
}

func (u URL) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

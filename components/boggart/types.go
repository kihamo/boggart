package boggart

import (
	"net"
	"net/url"
	"time"
)

type IP struct {
	net.IP
}

func (i IP) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

type IPMask struct {
	net.IPMask
}

func (m *IPMask) UnmarshalText(text []byte) error {
	ip := &net.IP{}

	err := ip.UnmarshalText(text)
	if err == nil {
		*m = IPMask{
			IPMask: net.IPMask(ip.To4()),
		}
	}

	return err
}

func (m IPMask) String() string {
	mask := net.CIDRMask(m.Size())
	return net.IP(mask).String()
}

type HardwareAddr struct {
	net.HardwareAddr
}

func (a HardwareAddr) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *HardwareAddr) UnmarshalText(text []byte) error {
	mac, err := net.ParseMAC(string(text))

	if err == nil {
		*a = HardwareAddr{
			HardwareAddr: mac,
		}
	}

	return err
}

type URL struct {
	url.URL
}

func (u URL) MarshalText() ([]byte, error) {
	return u.MarshalBinary()
}

type Location struct {
	time.Location
}

func (l Location) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Location) UnmarshalText(text []byte) error {
	location, err := time.LoadLocation(string(text))

	if err == nil {
		*l = Location{
			Location: *location,
		}
	}

	return err
}

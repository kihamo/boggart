package v3

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/serial"
	"github.com/kihamo/boggart/components/boggart/protocols/serial_network"
	m "github.com/kihamo/boggart/components/boggart/providers/mercury"
	mercury "github.com/kihamo/boggart/components/boggart/providers/mercury/v3"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	u, err := url.Parse(config.RS485Address)
	if err != nil {
		return nil, err
	}

	var conn m.Connection

	switch u.Scheme {
	case "tcp", "tcp4", "tcp6":
		conn = serial_network.NewTCPClient(u.Scheme, u.Host)

	case "udp", "udp4", "udp6", "unixgram":
		conn = serial_network.NewUDPClient(u.Scheme, u.Host)

	default:
		conn = serial.Dial(config.RS485Address, serial.WithTimeout(config.RS485Timeout))
	}

	opts := make([]mercury.Option, 0)
	if config.Address != "" {
		opts = append(opts, mercury.WithAddress(mercury.ConvertSerialNumber(config.Address)))
	}

	return &Bind{
		provider: mercury.New(conn, opts...),
		config:   config,

		tariff1:   atomic.NewUint32Null(),
		voltage1:  atomic.NewFloat32Null(),
		voltage2:  atomic.NewFloat32Null(),
		voltage3:  atomic.NewFloat32Null(),
		amperage1: atomic.NewFloat32Null(),
		amperage2: atomic.NewFloat32Null(),
		amperage3: atomic.NewFloat32Null(),
		power1:    atomic.NewFloat32Null(),
		power2:    atomic.NewFloat32Null(),
		power3:    atomic.NewFloat32Null(),
	}, nil
}

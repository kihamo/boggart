package v3

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	"github.com/kihamo/boggart/protocols/serial"
	"github.com/kihamo/boggart/protocols/serial_network"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	u, err := url.Parse(config.RS485Address)
	if err != nil {
		return nil, err
	}

	var conn connection.Conn

	switch u.Scheme {
	case "tcp", "tcp4", "tcp6":
		conn = serial_network.NewTCPClient(u.Scheme, u.Host)

	case "udp", "udp4", "udp6", "unixgram":
		conn = serial_network.NewUDPClient(u.Scheme, u.Host)

	default:
		conn = connection.NewIO(
			serial.Dial(serial.WithTarget(config.RS485Address),
				serial.WithTimeout(config.RS485Timeout),
			))
	}

	opts := make([]mercury.Option, 0)
	if config.Address != "" {
		opts = append(opts, mercury.WithAddressAsString(config.Address))
	}

	return &Bind{
		provider: mercury.New(conn, opts...),
		config:   config,
	}, nil
}

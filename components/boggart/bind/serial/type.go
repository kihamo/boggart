package serial

import (
	"errors"
	"net"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/serial"
	"github.com/kihamo/boggart/protocols/serial_network"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{}

	dial := serial.Dial(config.Target, serial.WithTimeout(config.Timeout))
	address := net.JoinHostPort(config.Host, strconv.FormatInt(config.Port, 10))

	switch config.Network {
	case "tcp", "tcp4", "tcp6":
		bind.server = serial_network.NewTCPServer(config.Network, address, dial)

	case "udp", "udp4", "udp6":
		bind.server = serial_network.NewUDPServer(config.Network, address, dial)

	default:
		return nil, errors.New("unsupported network " + config.Network)
	}

	return bind, nil
}

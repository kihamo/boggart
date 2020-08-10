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

	opts := []serial.Option{
		serial.WithAddress(config.Target),
		serial.WithBaudRate(config.BaudRate),
		serial.WithDataBits(config.DataBits),
		serial.WithStopBits(config.StopBits),
		serial.WithParity(config.Parity),
		serial.WithTimeout(config.Timeout),
		serial.WithOnce(config.Once),
	}

	dial := serial.Dial(opts...)
	address := net.JoinHostPort(config.Host, strconv.FormatInt(config.Port, 10))

	switch config.Network {
	case "tcp", "tcp4", "tcp6":
		bind.server = serialnetwork.NewTCPServer(config.Network, address, dial)

	case "udp", "udp4", "udp6":
		bind.server = serialnetwork.NewUDPServer(config.Network, address, dial)

	default:
		return nil, errors.New("unsupported network " + config.Network)
	}

	return bind, nil
}

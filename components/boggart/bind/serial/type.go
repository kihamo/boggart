package serial

import (
	"net"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/serial"
	"github.com/kihamo/boggart/components/boggart/protocols/serial_tcp"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	dial := serial.Dial(config.Target, serial.WithTimeout(config.Timeout))
	address := net.JoinHostPort(config.Host, strconv.FormatInt(config.Port, 10))

	return &Bind{
		server: serial_tcp.NewServer(address, dial),
	}, nil
}

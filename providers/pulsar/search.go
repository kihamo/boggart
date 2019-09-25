package pulsar

import (
	"github.com/kihamo/boggart/protocols/connection"
)

var commandSearch = []byte{0xF0, 0x0F, 0x0F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5, 0x44}

func DeviceAddress(conn connection.Conn) ([]byte, error) {
	response, err := connection.NewInvoker(conn).Invoke(commandSearch)
	if err != nil {
		return nil, err
	}

	return response[4:8], nil
}

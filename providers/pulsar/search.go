package pulsar

import (
	"github.com/kihamo/boggart/protocols/serial"
)

func DeviceAddress(c *serial.Serial) ([]byte, error) {
	response, err := c.Invoke([]byte{0xF0, 0x0F, 0x0F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5, 0x44})
	if err != nil {
		return nil, err
	}

	return response[4:8], nil
}

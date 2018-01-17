package pulsar

import (
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
)

func DeviceAddress(c *rs485.Connection) ([]byte, error) {
	response, err := c.Request([]byte{0xF0, 0x0F, 0x0F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5, 0x44})
	if err != nil {
		return nil, err
	}

	return response[4:8], nil
}

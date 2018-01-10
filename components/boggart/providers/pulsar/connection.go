package pulsar

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/goburrow/serial"
)

const (
	DefaultSerialAddress = "/dev/ttyUSB0"
	DefaultTimeout       = time.Second
)

type Connection struct {
	config *serial.Config
}

func NewConnection(address string, timeout time.Duration) *Connection {
	if address == "" {
		address = DefaultSerialAddress
	}

	if timeout < 0 {
		timeout = DefaultTimeout
	}

	return &Connection{
		config: &serial.Config{
			BaudRate: 9600,
			Parity:   "N",
			Address:  address,
			Timeout:  timeout,
		},
	}
}

func (c *Connection) Address() string {
	return c.config.Address
}

func (c *Connection) Timeout() time.Duration {
	return c.config.Timeout
}

func (c *Connection) DeviceAddress() ([]byte, error) {
	response, err := c.RequestRaw([]byte{0xF0, 0x0F, 0x0F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5, 0x44})
	if err != nil {
		return nil, err
	}

	return response[4:8], nil
}

func (c *Connection) RequestRaw(request []byte) ([]byte, error) {
	port, err := serial.Open(c.config)

	if err != nil {
		return nil, err
	}
	defer port.Close()

	if _, err := port.Write(request); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)

	for {
		b := make([]byte, 512)
		n, err := port.Read(b)
		if err != nil {
			break
		}

		if n != 0 {
			buffer.Write(b[:n])
		}
	}

	return buffer.Bytes(), err
}

func (c *Connection) Request(address []byte, function byte, data []byte) ([]byte, error) {
	var request []byte

	// device address
	request = append(request, address...)

	// function
	request = append(request, function)

	// length of packet
	l := len(request) + 1 + len(data) + 2 + 2
	request = append(request, byte(l))

	// data in
	request = append(request, data...)

	// request id
	requestId := GenerateRequestId()
	request = append(request, requestId...)

	// check sum CRC16
	request = append(request, GenerateCRC16(request)...)

	// fmt.Println("Request: ", request, p.ToString(request))

	response, err := c.RequestRaw(request)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Response: ", response, p.ToString(response))

	l = len(response)
	if l < 10 {
		return nil, errors.New("Error length of response packet")
	}

	// check crc16
	crc16 := GenerateCRC16(response[:l-2])
	if bytes.Compare(response[l-2:], crc16) != 0 {
		return nil, errors.New("Error CRC16 of response packet")
	}

	// check id
	if bytes.Compare(response[l-(2+len(requestId)):l-2], requestId) != 0 {
		return nil, errors.New("Error ID of response packet")
	}

	// check error
	if response[4] == FunctionBadCommand {
		return nil, fmt.Errorf("Device returns error code #%d", response[6])
	}

	return response[6 : l-4], nil
}

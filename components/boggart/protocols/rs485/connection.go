package rs485

import (
	"bytes"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

const (
	DefaultSerialAddress = "/dev/ttyUSB0"
	DefaultTimeout       = time.Second
)

type Connection struct {
	lock sync.Mutex

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

func (c *Connection) Request(request []byte) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

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

package serial

import (
	"bytes"
	"sync"

	"github.com/goburrow/serial"
)

var (
	multiRequestsMutex       sync.RWMutex
	multiRequestsConnections = make(map[string]sync.RWMutex)
)

type Connection struct {
	target  string
	options dialOptions
}

func Dial(target string, opts ...DialOption) *Connection {
	if target == "" {
		target = DefaultSerialAddress
	}

	conn := &Connection{
		target:  target,
		options: defaultDialOptions(),
	}

	for _, opt := range opts {
		opt.apply(&conn.options)
	}

	if !conn.options.allowMultiRequest {
		multiRequestsMutex.Lock()
		multiRequestsConnections[target] = sync.RWMutex{}
		multiRequestsMutex.Unlock()
	}

	return conn
}

func (c *Connection) Request(request []byte) ([]byte, error) {
	multiRequestsMutex.Lock()
	lock, ok := multiRequestsConnections[c.target]
	multiRequestsMutex.Unlock()

	if ok {
		lock.Lock()
		defer lock.Unlock()
	}

	port, err := serial.Open(&serial.Config{
		BaudRate: 9600,
		Parity:   "N",
		Address:  c.target,
		Timeout:  c.options.timeout,
	})

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

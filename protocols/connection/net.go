package connection

import (
	"errors"
	"net"
	"time"
)

type Net struct {
	address string
	options options
}

func Dial(address string, opts ...Option) *Net {
	conn := &Net{
		address: address,
		options: options{
			network: "tcp",
		},
	}

	for _, opt := range opts {
		opt.apply(&conn.options)
	}

	return conn
}

func (c *Net) connect() (net.Conn, error) {
	conn, err := net.Dial(c.options.network, c.address)
	if err != nil {
		return nil, err
	}

	tcp, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, errors.New("failed cast connect to *net.TCPConn")
	}

	return tcp, err
}

func (c *Net) Read(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	if c.options.readTimeout > 0 {
		if err = conn.SetReadDeadline(time.Now().Add(c.options.readTimeout)); err != nil {
			return -1., err
		}
	}

	return conn.Read(b)
}

func (c *Net) Write(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	if c.options.writeTimeout > 0 {
		if err = conn.SetWriteDeadline(time.Now().Add(c.options.readTimeout)); err != nil {
			return -1., err
		}
	}

	return conn.Write(b)
}

func (c *Net) Close() error {
	return nil
}

func (c *Net) Invoke(request []byte) (response []byte, err error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(request); err != nil {
		return nil, err
	}

	b := make([]byte, bufferSize)

	n, err := conn.Read(b)
	if n > 0 {
		return b[:n], nil
	}

	return nil, err
}

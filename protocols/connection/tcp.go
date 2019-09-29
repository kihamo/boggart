package connection

import (
	"errors"
	"net"
	"time"
)

type TCP struct {
	address string
	options options
}

func DialTCP(address string, opts ...Option) *TCP {
	conn := &TCP{
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

func (c *TCP) connect() (*net.TCPConn, error) {
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

func (c *TCP) Read(b []byte) (n int, err error) {
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

func (c *TCP) Write(b []byte) (n int, err error) {
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

func (c *TCP) Close() error {
	return nil
}

func (c *TCP) Invoke(request []byte) (response []byte, err error) {
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

	return nil, nil
}

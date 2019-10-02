package connection

import (
	"net"
	"sync"
	"time"
)

type Net struct {
	address string
	options options
	once    sync.Once
	conn    net.Conn
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

func (c *Net) connect() (conn net.Conn, err error) {
	if c.options.once {
		c.once.Do(func() {
			c.conn, err = net.Dial(c.options.network, c.address)
		})

		conn = c.conn
	} else {
		conn, err = net.Dial(c.options.network, c.address)
	}

	return conn, err
}

func (c *Net) Read(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	if c.options.readTimeout > 0 {
		if err = conn.SetReadDeadline(time.Now().Add(c.options.readTimeout)); err != nil {
			return -1, err
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
		if err = conn.SetWriteDeadline(time.Now().Add(c.options.writeTimeout)); err != nil {
			return -1, err
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

	if c.options.writeTimeout > 0 {
		if err = conn.SetWriteDeadline(time.Now().Add(c.options.writeTimeout)); err != nil {
			return nil, err
		}
	}

	if _, err = conn.Write(request); err != nil {
		return nil, err
	}

	b := make([]byte, bufferSize)

	if c.options.readTimeout > 0 {
		if err = conn.SetReadDeadline(time.Now().Add(c.options.readTimeout)); err != nil {
			return nil, err
		}
	}

	n, err := conn.Read(b)
	if n > 0 {
		return b[:n], err
	}

	return nil, err
}

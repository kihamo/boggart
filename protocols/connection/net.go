package connection

import (
	"net"
	"sync"
)

var readBufferPool sync.Pool

func init() {
	readBufferPool.New = func() interface{} {
		buf := make([]byte, bufferSize)
		return &buf
	}
}

type Net struct {
	address string
	options options
	once    *sync.Once
	conn    net.Conn
}

func Dial(address string, opts ...Option) *Net {
	conn := &Net{
		address: address,
		options: options{
			network: "tcp",
		},
		once: new(sync.Once),
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
			if err != nil {
				c.once = new(sync.Once)
			}
		})

		conn = c.conn
	} else {
		conn, err = net.Dial(c.options.network, c.address)
	}

	return conn, err
}

func (c *Net) Read(b []byte) (int, error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	if err = SetDeadline(c.options.readTimeout, conn.SetReadDeadline); err != nil {
		return -1, err
	}

	return conn.Read(b)
}

func (c *Net) Write(b []byte) (int, error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	if err = SetDeadline(c.options.writeTimeout, conn.SetWriteDeadline); err != nil {
		return -1, err
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

	if err = SetDeadline(c.options.writeTimeout, conn.SetWriteDeadline); err != nil {
		return nil, err
	}

	if _, err = conn.Write(request); err != nil {
		return nil, err
	}

	if err = SetDeadline(c.options.readTimeout, conn.SetReadDeadline); err != nil {
		return nil, err
	}

	buf := readBufferPool.Get().(*[]byte)
	defer readBufferPool.Put(buf)

	n, err := conn.Read(*buf)
	if n > 0 {
		return append([]byte(nil), (*buf)[:n]...), err
	}

	return nil, err
}

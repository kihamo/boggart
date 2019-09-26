package connection

import (
	"io"
	"sync"
)

type IO struct {
	conn io.ReadWriter
}

func NewIO(rw io.ReadWriter) Conn {
	return &IO{
		conn: rw,
	}
}

func (c *IO) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *IO) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *IO) Close() error {
	if closer, ok := c.conn.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func (c *IO) Lock() {
	if locker, ok := c.conn.(sync.Locker); ok {
		locker.Lock()
	}
}

func (c *IO) Unlock() {
	if locker, ok := c.conn.(sync.Locker); ok {
		locker.Unlock()
	}
}

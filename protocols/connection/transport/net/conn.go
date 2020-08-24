package net

import (
	"errors"
	"net"
	"time"

	"github.com/kihamo/boggart/protocols/connection/transport"
)

func SetDeadline(duration time.Duration, f func(t time.Time) error) error {
	var deadline time.Time

	if duration > 0 {
		deadline = time.Now().Add(duration)
	} else {
		deadline = time.Time{}
	}

	return f(deadline)
}

type Net struct {
	options Options
	conn    net.Conn
}

func New(options ...Option) *Net {
	c := &Net{
		options: DefaultOptions(),
	}

	for _, option := range options {
		option.apply(&c.options)
	}

	return c
}

func (c *Net) Dial() (_ transport.Transport, err error) {
	w := &Net{
		options: c.options,
	}

	w.conn, err = net.Dial(c.options.network, c.options.address)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (c *Net) Read(b []byte) (_ int, err error) {
	if c.conn == nil {
		return -1, errors.New("connection isn't init")
	}

	if err = SetDeadline(c.options.readTimeout, c.conn.SetReadDeadline); err != nil {
		return -1, err
	}

	return c.conn.Read(b)
}

func (c *Net) Write(b []byte) (_ int, err error) {
	if c.conn == nil {
		return -1, errors.New("connection isn't init")
	}

	if err = SetDeadline(c.options.writeTimeout, c.conn.SetWriteDeadline); err != nil {
		return -1, err
	}

	return c.conn.Write(b)
}

func (c *Net) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *Net) Options() map[string]interface{} {
	return c.options.Map()
}

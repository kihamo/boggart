package internal

import (
	"bytes"
	"context"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio/internal/packet"
)

const (
	DefaultPort = 54321
	network     = "udp"
)

type Connection struct {
	io.Closer

	conn *net.UDPConn
}

func NewConnection(address string) (*Connection, error) {
	if _, _, err := net.SplitHostPort(address); err != nil {
		address = address + ":" + strconv.Itoa(DefaultPort)
	}

	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(network, nil, addr)
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn: conn,
	}, nil
}

func (c *Connection) Invoke(ctx context.Context, request io.WriterTo, response io.ReaderFrom) (err error) {
	done := make(chan error, 1)

	go func() {
		var deadline time.Time

		if d, ok := ctx.Deadline(); ok {
			deadline = d
		} else {
			// disable deadline
			deadline = time.Time{}
		}

		if err := c.conn.SetDeadline(deadline); err != nil {
			done <- err
			return
		}

		if _, err := request.WriteTo(c.conn); err != nil {
			done <- err
			return
		}

		if response != nil {
			b := make([]byte, packet.MaxBufferSize)
			n, _, err := c.conn.ReadFromUDP(b)

			if err != nil {
				done <- err
				return
			}

			if n > 0 {
				_, err = response.ReadFrom(bytes.NewBuffer(b[:n]))
				if err != nil {
					done <- err
					return
				}
			}
		}

		done <- nil
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()

	case err = <-done:
	}

	return err
}

func (c *Connection) LocalAddr() *net.UDPAddr {
	return c.conn.LocalAddr().(*net.UDPAddr)
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

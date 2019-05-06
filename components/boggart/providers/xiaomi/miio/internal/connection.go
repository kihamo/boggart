package internal

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/internal/packet"
)

const (
	DefaultPort     = 54321
	DefaultDeadline = time.Second * 5

	network = "udp"
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

func (c *Connection) Invoke(request packet.Packet, response packet.Packet) error {
	if _, err := request.WriteTo(c.conn); err != nil {
		return err
	}

	if response == nil {
		return nil
	}

	err := c.conn.SetDeadline(time.Now().Add(DefaultDeadline))
	if err != nil {
		return err
	}

	b := make([]byte, packet.MaxBufferSize)
	n, _, err := c.conn.ReadFromUDP(b)

	if n > 0 {
		_, err = response.ReadFrom(bytes.NewBuffer(b[:n]))
		if err != nil {
			return err
		}
	}

	return err
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

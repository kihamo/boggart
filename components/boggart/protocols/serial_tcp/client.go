package serial_tcp

import (
	"bytes"
	"net"
	"time"
)

type Client struct {
	net.Conn

	conn *net.TCPConn
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	tcp := conn.(*net.TCPConn)

	err = tcp.SetKeepAlive(true)
	if err != nil {
		return nil, err
	}

	err = tcp.SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: tcp,
	}, nil
}

func (c *Client) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func (c *Client) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Client) Invoke(request []byte) (response []byte, err error) {
	if _, err = c.conn.Write(request); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)

	for {
		b := make([]byte, 512)
		n, e := c.conn.Read(b)
		if e != nil {
			break
		}

		if n != 0 {
			buffer.Write(b[:n])
		}
	}

	return buffer.Bytes(), nil
}

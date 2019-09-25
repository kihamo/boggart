package serial_network

import (
	"io/ioutil"
	"net"
)

type TCPClient struct {
	network string
	address string
}

func NewTCPClient(network, address string) *TCPClient {
	return &TCPClient{
		network: network,
		address: address,
	}
}

func (c *TCPClient) connect() (net.Conn, error) {
	return net.Dial(c.network, c.address)
}

func (c *TCPClient) Read(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	return conn.Read(b)
}

func (c *TCPClient) Write(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	return conn.Write(b)
}

func (c *TCPClient) Close() error {
	return nil
}

func (c *TCPClient) Invoke(request []byte) (response []byte, err error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(request); err != nil {
		return nil, err
	}

	return ioutil.ReadAll(conn)
}

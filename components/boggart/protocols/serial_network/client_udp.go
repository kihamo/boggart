package serial_network

import (
	"errors"
	"net"
	"time"
)

type UDPClient struct {
	network string
	address string
}

func NewUDPClient(network, address string) *UDPClient {
	return &UDPClient{
		network: network,
		address: address,
	}
}

func (c *UDPClient) connect() (*net.UDPConn, error) {
	conn, err := net.Dial(c.network, c.address)
	if err != nil {
		return nil, err
	}

	udp, ok := conn.(*net.UDPConn)
	if !ok {
		return nil, errors.New("failed cast connect to *net.UDPConn")
	}

	/*
		err = udp.SetDeadline(time.Now().Add(time.Second*2))
		if err != nil {
			return nil, err
		}
	*/

	//udp.SetReadDeadline(time.Now().Add(time.Second * 2))

	return udp, err
}

func (c *UDPClient) Read(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	return conn.Read(b)
}

func (c *UDPClient) Write(b []byte) (n int, err error) {
	conn, err := c.connect()
	if err != nil {
		return -1, err
	}

	return conn.Write(b)
}

func (c *UDPClient) Invoke(request []byte) (response []byte, err error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(request); err != nil {
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		return nil, err
	}

	b := make([]byte, maxBufferSize)

	n, _, err := conn.ReadFromUDP(b)
	if n > 0 {
		return b[:n], nil
	}

	return nil, nil
}

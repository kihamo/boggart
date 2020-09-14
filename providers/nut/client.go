package nut

import (
	"net"
	"strconv"
)

const (
	DefaultPort = 3493
)

type Client struct {
	dsn      string
	username string
	password string
}

func New(host, username, password string) *Client {
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = host + ":" + strconv.Itoa(DefaultPort)
	}

	return &Client{
		dsn:      "tcp://" + host + "?read-timeout=10s&write-timeout=10s&once=true",
		username: username,
		password: password,
	}
}

func (c *Client) Session() *Session {
	return NewSession(c.dsn, c.username, c.password)
}

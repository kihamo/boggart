package mikrotik

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/routeros.v2"
)

type Client struct {
	client *routeros.Client
}

func NewClient(address, username, password string, timeout time.Duration) (*Client, error) {
	client, err := routeros.DialTimeout(address, username, password, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) System() (map[string]string, error) {
	reply, err := c.client.RunArgs([]string{"/system/resource/print"})
	if err != nil {
		return nil, nil
	}

	if len(reply.Re) == 0 {
		return nil, errors.New("Empty reply from device")
	}

	return reply.Re[0].Map, nil
}

func (c *Client) WifiClients() ([]interface{}, error) {
	reply, err := c.client.RunArgs([]string{"/interface/wireless/registration-table/print"})
	if err != nil {
		return nil, err
	}

	for i, re := range reply.Re {
		fmt.Println(i, re)
	}

	return nil, nil
}

package modbus

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/goburrow/modbus"
)

/*
FuncCodeReadHoldingRegisters       = 3
FuncCodeWriteSingleRegister        = 6
*/

const (
	DefaultPort = 503
)

type Client struct {
	modbus.Client

	handler modbus.ClientHandler
	options options
}

func NewClient(address *url.URL, opts ...Option) *Client {
	client := &Client{
		options: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&client.options)
	}

	switch address.Scheme {
	case "rtu":
		if address.Port() == "" {
			address.Host = net.JoinHostPort(address.Host, strconv.Itoa(DefaultPort))
		}

		handler := modbus.NewRTUClientHandler(address.Host)
		handler.SlaveId = client.options.slaveID
		handler.Timeout = client.options.timeout
		handler.IdleTimeout = client.options.idleTimeout

		if client.options.logger != nil {
			handler.Logger = log.New(client.options.logger, "", 0)
		}

		client.handler = handler

	default:
		if address.Port() == "" {
			address.Host = net.JoinHostPort(address.Host, strconv.Itoa(DefaultPort))
		}

		handler := modbus.NewTCPClientHandler(address.Host)
		handler.SlaveId = client.options.slaveID
		handler.Timeout = client.options.timeout
		handler.IdleTimeout = client.options.idleTimeout

		if client.options.logger != nil {
			handler.Logger = log.New(client.options.logger, "", 0)
		}

		client.handler = handler
	}

	client.Client = modbus.NewClient(client.handler)

	return client
}

func (c *Client) Close() error {
	if closer, ok := c.handler.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func (c *Client) Read(address, quantity uint16) (response []byte, err error) {
	for trie := uint8(1); trie <= c.options.maxTries; trie++ {
		response, err = c.ReadHoldingRegisters(address, quantity)
		if err == nil {
			break
		}
	}

	return response, err
}

func (c *Client) ReadUint16(address uint16) (value uint16, err error) {
	response, err := c.Read(address, 1)
	if err == nil {
		return binary.BigEndian.Uint16(response), err
	}

	return value, err
}

func (c *Client) Write(address, quantity uint16, payload []byte) (err error) {
	var response []byte

	for trie := uint8(1); trie <= c.options.maxTries; trie++ {
		response, err = c.WriteMultipleRegisters(address, quantity, payload)

		if err == nil {
			_ = binary.BigEndian.Uint16(response)

			//if code == writeResponseSuccess {
			//	break
			//}
			//
			//err = fmt.Errorf("device return not success response %d", code)
		}
	}

	return err
}

func (c *Client) WriteUint16(address, value uint16) error {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, value)

	return c.Write(address, 2, payload)
}

func (c *Client) WriteUint32(address uint16, value uint32) error {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, value)

	return c.Write(address, 2, payload)
}

package modbus

import (
	"encoding/binary"
	"fmt"
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

func (c *Client) CallWithTriesLimit(f func() ([]byte, error)) (response []byte, err error) {
	for trie := uint8(1); trie <= c.options.maxTries; trie++ {
		response, err = f()
		if err == nil {
			break
		}
	}

	return response, err
}

func (c *Client) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	return c.CallWithTriesLimit(func() ([]byte, error) {
		return c.Client.ReadHoldingRegisters(address, quantity)
	})
}

func (c *Client) ReadHoldingRegistersAsMap(address, quantity uint16) (result map[uint16]uint16, err error) {
	response, err := c.ReadHoldingRegisters(address, quantity)

	if err != nil {
		return nil, err
	}

	if len(response) != int(quantity)*2 {
		return nil, fmt.Errorf("wrong response payload length %d need %d", len(response), quantity*2)
	}

	result = make(map[uint16]uint16, int(quantity))

	for i := uint16(0); i < quantity; i++ {
		result[address+i] = binary.BigEndian.Uint16(response[i*2 : i*2+2])
	}

	return result, err
}

func (c *Client) ReadHoldingRegistersUint8(address uint16) (value uint8, err error) {
	response, err := c.ReadHoldingRegisters(address, 1)
	if err == nil {
		return response[0], err
	}

	return value, err
}

func (c *Client) ReadHoldingRegistersUint16(address uint16) (value uint16, err error) {
	response, err := c.ReadHoldingRegisters(address, 1)
	if err == nil {
		return binary.BigEndian.Uint16(response), err
	}

	return value, err
}

func (c *Client) ReadHoldingRegistersBool(address uint16) (value bool, err error) {
	response, err := c.ReadHoldingRegisters(address, 1)
	if err == nil {
		return response[1] == 1, err
	}

	return value, err
}

func (c *Client) WriteSingleRegister(address, payload uint16) ([]byte, error) {
	return c.CallWithTriesLimit(func() ([]byte, error) {
		return c.Client.WriteSingleRegister(address, payload)
	})
}

func (c *Client) WriteSingleRegisterUint16Bytes(address uint16, payload []byte) ([]byte, error) {
	return c.WriteSingleRegister(address, binary.BigEndian.Uint16(payload))
}

func (c *Client) WriteMultipleRegisters(address, quantity uint16, payload []byte) ([]byte, error) {
	return c.CallWithTriesLimit(func() ([]byte, error) {
		return c.Client.WriteMultipleRegisters(address, quantity, payload)
	})
}

package mc6

import (
	"io"
	"log"
	"net/url"

	"github.com/goburrow/modbus"
)

const (
	AddressRoomTemperature     uint16 = 0
	AddressFloorTemperature    uint16 = 1
	AddressHumidity            uint16 = 2
	AddressHeatingOutputStatus uint16 = 9
	AddressDeviceType          uint16 = 18

	DeviceTypeHotWater             uint16 = 2
	DeviceTypeElectricHeating      uint16 = 3
	DeviceTypeFCU2                 uint16 = 4
	DeviceTypeFCU4                 uint16 = 5
	DeviceTypeBase                 uint16 = 30 // базовый простой MC6-HA, без горячей воды
	DeviceTypeElectricHeatingTimer uint16 = 31
)

type MC6 struct {
	handler    modbus.ClientHandler
	connection modbus.Client
	options    options
}

func New(address *url.URL, opts ...Option) *MC6 {
	m := &MC6{
		options: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&m.options)
	}

	switch address.Scheme {
	default:
		handler := modbus.NewTCPClientHandler(address.Host)
		handler.SlaveId = m.options.slaveID
		handler.Timeout = m.options.timeout
		handler.IdleTimeout = m.options.idleTimeout

		if m.options.logger != nil {
			handler.Logger = log.New(m.options.logger, "", 0)
		}

		m.handler = handler
	}

	m.connection = modbus.NewClient(m.handler)

	return m
}

func (m *MC6) Close() error {
	if closer, ok := m.handler.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func (m *MC6) Invoke(address uint16) (response []byte, err error) {
	for trie := 1; trie <= m.options.maxTries; trie++ {
		response, err = m.connection.ReadHoldingRegisters(address, 1)
		if err == nil {
			return response, err
		}
	}

	return response, err
}

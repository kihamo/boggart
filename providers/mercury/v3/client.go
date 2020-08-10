package v3

import (
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/connection"
)

type MercuryV3 struct {
	connection connection.Connection
	options    options
}

func New(conn connection.Connection, opts ...Option) *MercuryV3 {
	conn.ApplyOptions(connection.WithLock(true))
	conn.ApplyOptions(connection.WithOnceInit(true))

	m := &MercuryV3{
		connection: conn,
		options:    defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&m.options)
	}

	return m
}

func (m *MercuryV3) Invoke(request *Request) (response *Response, err error) {
	err = m.ChannelOpen(m.options.accessLevel, m.options.password)
	if err != nil {
		return response, err
	}

	response, err = m.InvokeRaw(request)
	if err == nil {
		err = m.ChannelClose()
	}

	return response, err
}

func (m *MercuryV3) InvokeRaw(request *Request) (*Response, error) {
	if request.Address() == 0 {
		request = request.WithAddress(m.options.address)
	}

	if request.address >= AddressReservedBegin {
		return nil, errors.New("device address can't be reserved 0xF1...0xFF")
	}

	requestData, err := request.MarshalBinary()
	if err != nil {
		return nil, err
	}

	responseData, err := m.connection.Invoke(requestData)
	if err != nil {
		return nil, err
	}

	response := NewResponse()
	if err = response.UnmarshalBinary(responseData); err != nil {
		return nil, err
	}

	return response, err
}

func (m *MercuryV3) Raw() error {
	bwri := uint8(0x8<<4) | uint8(PhaseNumber1)
	resp, err := m.AuxiliaryParameters(bwri)

	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}

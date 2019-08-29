package v3

import (
	"fmt"

	"github.com/kihamo/boggart/components/boggart/providers/mercury"
)

type MercuryV3 struct {
	connection mercury.Connection
	options    options
}

func New(connection mercury.Connection, opts ...Option) *MercuryV3 {
	m := &MercuryV3{
		connection: connection,
		options:    defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&m.options)
	}

	return m
}

func (m *MercuryV3) Request(request *Request) (response *Response, err error) {
	err = m.ChannelOpen(m.options.accessLevel, m.options.password)
	if err != nil {
		return response, err
	}

	response, err = m.RequestRaw(request)
	if err == nil {
		err = m.ChannelClose()
	}

	return response, err
}

func (m *MercuryV3) RequestRaw(request *Request) (*Response, error) {
	// fmt.Println("Request: >>>>>")
	// fmt.Println(hex.Dump(request.Bytes()))

	data, err := m.connection.Invoke(request.Bytes())
	if err != nil {
		return nil, err
	}

	return ParseResponse(data)

	/*
		response, err := ParseResponse(data)
		if err == nil {
			fmt.Println("Response: <<<<<")
			fmt.Println(hex.Dump(response.Bytes()))
		}

		return response, err
	*/
}

func (m *MercuryV3) Raw() error {
	bwri := int64(0x8<<4) | int64(PhaseNumber1)
	resp, err := m.AuxiliaryParameters(bwri)

	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}

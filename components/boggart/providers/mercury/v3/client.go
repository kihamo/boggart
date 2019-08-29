package v3

import (
	"encoding/hex"
	"fmt"

	"github.com/kihamo/boggart/components/boggart/providers/mercury"
)

type MercuryV3 struct {
	address    byte
	connection mercury.Connection
}

func New(connection mercury.Connection) *MercuryV3 {
	return &MercuryV3{
		address:    0x0,
		connection: connection,
	}
}

func (m *MercuryV3) WithAddress(address byte) *MercuryV3 {
	m.address = address
	return m
}

func (m *MercuryV3) Request(request *Request) (*Response, error) {
	fmt.Println("Request: >>>>>")
	fmt.Println(hex.Dump(request.Bytes()))

	data, err := m.connection.Invoke(request.Bytes())
	if err != nil {
		return nil, err
	}

	response, err := ParseResponse(data)
	if err == nil {
		fmt.Println("Response: <<<<<")
		fmt.Println(hex.Dump(response.Bytes()))
	}

	return response, err
}

func (m *MercuryV3) Raw() error {
	resp, err := m.Request(&Request{
		Address:            m.address,
		Code:               RequestCodeReadParameter,
		ParameterCode:      &[]byte{ParamCodeAuxiliaryParameters3}[0],
		ParameterExtension: &[]byte{0x12}[0],
	})

	if err != nil {
		return err
	}

	fmt.Println(resp.Payload)

	return nil
}

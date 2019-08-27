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

func New(address byte, connection mercury.Connection) *MercuryV3 {
	return &MercuryV3{
		address:    address,
		connection: connection,
	}
}

func (d *MercuryV3) Request(request *Request) (*Response, error) {
	fmt.Println("Request: >>>>>")
	fmt.Println(hex.Dump(request.Bytes()))

	data, err := d.connection.Invoke(request.Bytes())
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

// 2.1. ЗАПРОС НА ТЕСТИРОВАНИЕ КАНАЛА СВЯЗИ
func (d *MercuryV3) ChannelTest() error {
	resp, err := d.Request(&Request{
		Address: d.address,
		Code:    RequestCodeChannelTest,
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (d *MercuryV3) ChannelOpen(level accessLevel, password LevelPassword) error {
	l := byte(level)

	resp, err := d.Request(&Request{
		Address:       d.address,
		Code:          RequestCodeChannelOpen,
		ParameterCode: &l,
		Parameters:    password.Bytes(),
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (d *MercuryV3) ChannelClose() error {
	resp, err := d.Request(&Request{
		Address: d.address,
		Code:    RequestCodeChannelClose,
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}

func (d *MercuryV3) Raw() error {
	//b := byte(0x00)
	//tariff := byte(0xE0)

	resp, err := d.Request(&Request{
		Address: d.address,
		Code:    0x08,
		//ParameterCode: &b,
		Parameters: []byte{0x16, 0x2, 0x2},
	})

	if err != nil {
		return err
	}

	if err := ResponseError(resp); err != nil {
		return err
	}

	fmt.Println(resp.Payload)

	return nil
}

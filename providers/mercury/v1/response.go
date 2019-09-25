package v1

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/kihamo/boggart/protocols/serial"
)

type Response struct {
	Address []byte
	Command requestCommand
	Payload []byte
	CRC     []byte
}

func ParseResponse(data []byte) (*Response, error) {
	l := len(data)

	if l < 7 {
		return nil, errors.New("bad packet length")
	}

	response := &Response{
		Address: data[:4],
		Command: requestCommand(data[4]),
		Payload: data[5 : len(data)-2],
		CRC:     data[len(data)-2:],
	}

	crc := serial.GenerateCRC16(data[:l-2])
	if !bytes.Equal(crc, response.CRC) {
		return nil, errors.New("error CRC16 of response packet have " +
			hex.EncodeToString(crc) + " want " +
			hex.EncodeToString(response.CRC))
	}

	return response, nil
}

func (r *Response) Bytes() []byte {
	packet := append(r.Address, byte(r.Command))
	packet = append(packet, r.Payload...)
	packet = append(packet, r.CRC...)

	return packet
}

func (r *Response) String() string {
	return hex.EncodeToString(r.Bytes())
}

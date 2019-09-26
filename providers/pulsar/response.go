package pulsar

import (
	"encoding/hex"
	"errors"
)

type Response struct {
	Address   []byte
	Function  byte
	Length    byte
	ErrorCode byte
	Payload   []byte
	Id        []byte
	CRC       []byte
}

func ParseResponse(data []byte) (*Response, error) {
	l := len(data)

	if l < 10 {
		return nil, errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	r := &Response{
		Address:  data[:4],
		Function: data[4],
		Length:   data[5],
		Id:       data[l-4 : l-2],
		CRC:      data[l-2:],
	}

	if r.Function == FunctionBadCommand {
		r.ErrorCode = data[6]
	} else {
		r.Payload = data[6 : l-4]
	}

	return r, nil
}

func (r *Response) Bytes() []byte {
	packet := append(r.Address, r.Function)
	packet = append(packet, r.Length)
	packet = append(packet, r.Payload...)
	packet = append(packet, r.Id...)
	packet = append(packet, r.CRC...)

	return packet
}

func (r *Response) String() string {
	return hex.EncodeToString(r.Bytes())
}

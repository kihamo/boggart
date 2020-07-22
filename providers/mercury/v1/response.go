package v1

import (
	"bytes"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/kihamo/boggart/protocols/serial"
)

type Response struct {
	address []byte
	command requestCommand
	payload []byte
	crc     []byte

	lock sync.RWMutex
}

func ParseResponse(data []byte) (*Response, error) {
	l := len(data)

	if l < 7 {
		return nil, errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	response := &Response{
		address: data[:4],
		command: requestCommand(data[4]),
		payload: data[5 : l-2],
		crc:     data[l-2:],
	}

	crc := serial.GenerateCRC16(data[:l-2])
	if !bytes.Equal(crc, response.crc) {
		return nil, errors.New("error CRC16 of response packet " +
			hex.EncodeToString(data) + " have " +
			hex.EncodeToString(response.crc) + " want " +
			hex.EncodeToString(crc))
	}

	return response, nil
}

func (r *Response) Bytes() []byte {
	packet := append(r.Address(), byte(r.Command()))
	packet = append(packet, r.Payload()...)
	packet = append(packet, r.CRC()...)

	return packet
}

func (r *Response) Address() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.address...)
}

func (r *Response) Command() requestCommand {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.command
}

func (r *Response) Payload() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.payload...)
}

func (r *Response) PayloadAsBuffer() *Buffer {
	return NewBuffer(r.Payload())
}

func (r *Response) CRC() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.crc...)
}

func (r *Response) String() string {
	return hex.EncodeToString(r.Bytes())
}

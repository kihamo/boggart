package v1

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/kihamo/boggart/protocols/serial"
)

type Response struct {
	address uint32
	command uint8
	payload []byte
	crc     uint16

	lock sync.RWMutex
}

func ParseResponse(data []byte) (*Response, error) {
	l := len(data)

	if l < 7 {
		return nil, errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	response := &Response{
		address: binary.LittleEndian.Uint32(data[:4]),
		command: data[4],
		payload: data[5 : l-2],
		crc:     binary.LittleEndian.Uint16(data[l-2:]),
	}

	crc := serial.GenerateCRC16(data[:l-2])
	if !bytes.Equal(crc, data[l-2:]) {
		return nil, errors.New("error CRC16 of response packet " +
			hex.EncodeToString(data) + " have " +
			hex.EncodeToString(data[l-2:]) + " want " +
			hex.EncodeToString(crc))
	}

	return response, nil
}

func (r *Response) Bytes() []byte {
	var buf []byte

	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, r.Address())

	packet := append(buf, r.Command())
	packet = append(packet, r.Payload()...)

	buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, r.CRC())
	packet = append(packet, buf...)

	return packet
}

func (r *Response) Address() uint32 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.address
}

func (r *Response) Command() uint8 {
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

func (r *Response) CRC() uint16 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.crc
}

func (r *Response) String() string {
	return hex.EncodeToString(r.Bytes())
}

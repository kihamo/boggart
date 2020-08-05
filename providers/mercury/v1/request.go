package v1

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

type Request struct {
	address uint32
	command uint8
	payload []byte
}

func NewRequest(command uint8) *Request {
	return &Request{
		command: command,
	}
}

func (r *Request) Address() uint32 {
	return r.address
}

func (r *Request) WithAddress(address uint32) *Request {
	r.address = address
	return r
}

func (r *Request) WithPayload(payload []byte) *Request {
	r.payload = payload
	return r
}

func (r *Request) MarshalBinary() (_ []byte, err error) {
	packet := make([]byte, 4)
	binary.LittleEndian.PutUint32(packet, r.address)

	packet = append(packet, r.command)
	packet = append(packet, r.payload...)
	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet, nil
}

func (r *Request) String() string {
	data, _ := r.MarshalBinary()
	return hex.EncodeToString(data)
}

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

func (r *Request) WithAddress(address uint32) *Request {
	r.address = address
	return r
}

func (r *Request) WithPayload(payload []byte) *Request {
	r.payload = payload
	return r
}

func (r *Request) Bytes() []byte {
	packet := make([]byte, 4)
	binary.LittleEndian.PutUint32(packet, r.address)

	packet = append(packet, r.command)
	packet = append(packet, r.payload...)
	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

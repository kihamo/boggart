package pulsar

import (
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

type Request struct {
	Address  []byte
	Function byte
	Payload  []byte
	Id       []byte
}

func (r *Request) ID() []byte {
	if r.Id == nil {
		r.Id = serial.GenerateRequestID()
	}

	return r.Id
}

func (r *Request) Bytes() []byte {
	// device address + function
	packet := append(r.Address, r.Function)

	// length of packet
	l := len(packet) + 1 + len(r.Payload) + 2 + 2
	packet = append(packet, byte(l))

	// data in
	packet = append(packet, r.Payload...)

	// request id
	packet = append(packet, r.ID()...)

	// check sum CRC16
	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

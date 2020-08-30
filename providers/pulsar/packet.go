package pulsar

import (
	"encoding/hex"
	"errors"

	"github.com/kihamo/boggart/protocols/serial"
)

type Packet struct {
	address   []byte // 4 bytes
	function  uint8  // 1 byte
	length    uint8  // 1 byte
	errorCode ErrorCode
	payload   []byte
	id        []byte // 2 bytes
	crc       []byte // 2 bytes
}

func NewPacket() *Packet {
	return &Packet{}
}

func (p *Packet) Address() []byte {
	return append([]byte(nil), p.address...)
}

func (p *Packet) WithAddress(address []byte) *Packet {
	p.address = address
	return p
}

func (p *Packet) Length() uint8 {
	p.length = uint8(
		4 + // address
			1 + // function
			1 + // length
			len(p.Payload()) + // payload
			2 + // id
			2) // crc

	return p.length
}

func (p *Packet) Function() uint8 {
	return p.function
}

func (p *Packet) WithFunction(function uint8) *Packet {
	p.function = function
	return p
}

func (p *Packet) ErrorCode() ErrorCode {
	return p.errorCode
}

func (p *Packet) Payload() []byte {
	return append([]byte(nil), p.payload...)
}

func (p *Packet) WithPayload(payload []byte) *Packet {
	p.payload = payload
	return p
}

func (p *Packet) CRC() []byte {
	return append([]byte(nil), p.crc...)
}

func (p *Packet) ID() []byte {
	if p.id == nil {
		p.id = serial.GenerateRequestID()
	}

	return append([]byte(nil), p.id...)
}

func (p *Packet) MarshalBinary() (_ []byte, err error) {
	// device address + function
	packet := append(p.Address(), p.Function())

	// length of packet
	packet = append(packet, p.Length())

	// data
	packet = append(packet, p.Payload()...)

	// request id
	packet = append(packet, p.ID()...)

	// check sum CRC16
	p.crc = serial.GenerateCRC16(packet)
	packet = append(packet, p.crc...)

	return packet, nil
}

func (p *Packet) UnmarshalBinary(data []byte) (err error) {
	l := len(data)

	if l < 10 {
		return errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	p.address = data[:4]
	p.function = data[4]
	p.length = data[5]
	p.id = data[l-4 : l-2]
	p.crc = data[l-2:]

	if p.function == FunctionBadCommand {
		p.errorCode = ErrorCode(data[6])
	} else {
		p.payload = data[6 : l-4]
	}

	return nil
}

func (p *Packet) String() string {
	data, _ := p.MarshalBinary()
	return hex.EncodeToString(data)
}

func ReadCheck(data []byte) bool {
	if len(data) < 6 {
		return true
	}

	return uint8(len(data)) < data[5]
}

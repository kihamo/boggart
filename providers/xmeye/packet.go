package xmeye

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

var (
	regularPacketHeader = []byte{0xff, 0x01, 0x00, 0x00}
	payloadEOF          = []byte{0x0a, 0x00}
)

type packet struct {
	SessionID      uint32
	SequenceNumber uint32
	MessageID      uint16
	TotalPacket    uint16
	CurrentPacket  uint16
	PayloadLen     int
	Payload        *Payload
}

func (p packet) Marshal() []byte {
	message := make([]byte, 0x14) // build head

	// Head Flag, VERSION, RESERVED01, RESERVED02
	copy(message, regularPacketHeader)

	// SESSION ID
	binary.LittleEndian.PutUint32(message[0x04:], p.SessionID)

	// SEQUENCE NUMBER
	binary.LittleEndian.PutUint32(message[0x08:], p.SequenceNumber)

	// Total Packet
	if p.TotalPacket == 0 {
		p.TotalPacket = 1
	}

	message[0x0c] = byte(p.TotalPacket)

	// CurPacket
	if p.CurrentPacket == 0 {
		p.CurrentPacket = 1
	}

	message[0x0d] = byte(p.CurrentPacket)

	// Message Id
	binary.LittleEndian.PutUint16(message[0x0e:], p.MessageID)

	// Data Length
	var payloadLen uint32
	if p.Payload != nil {
		payloadLen = uint32(p.Payload.Len())
	}

	binary.LittleEndian.PutUint32(message[0x10:], payloadLen)

	if p.Payload != nil {
		// DATA
		message = append(message, p.Payload.Bytes()...)
		message = append(message, payloadEOF...)
	}

	return message
}

func (p packet) LoadPayload(payload interface{}) (err error) {
	if payload == nil {
		p.Payload.Reset()
		p.PayloadLen = 0
		return nil
	}

	if b, ok := payload.([]byte); ok {
		p.Payload.Reset()
		p.PayloadLen, err = p.Payload.Write(b)
		return err
	}

	encode, err := json.Marshal(payload)
	if err == nil {
		p.Payload.Reset()
		p.PayloadLen, err = p.Payload.Write(encode)
		return err
	}

	return err
}

func newPacket() *packet {
	return &packet{
		Payload: NewPayload(),
	}
}

func packetUnmarshal(message []byte) (*packet, error) {
	if !bytes.Equal(regularPacketHeader, message[:0x04]) {
		return nil, errors.New("invalid regular packet header")
	}

	packet := newPacket()

	packet.SessionID = binary.LittleEndian.Uint32(message[0x04:0x08])
	packet.SequenceNumber = binary.LittleEndian.Uint32(message[0x08:0x0c])
	packet.TotalPacket = uint16(message[0x0c])
	packet.CurrentPacket = uint16(message[0x0d])
	packet.MessageID = binary.LittleEndian.Uint16(message[0x0e:0x10])
	packet.PayloadLen = int(binary.LittleEndian.Uint32(message[0x10:0x14]))
	packet.Payload.Write(message[0x14:])

	return packet, nil
}

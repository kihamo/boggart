package xmeye

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

var (
	regularPacketHeader = []byte{0xff, 0x01, 0x00, 0x00}
	payloadEOF          = []byte{0x0a}
)

type Packet struct {
	sessionID      uint32
	sequenceNumber uint32
	messageID      uint16
	totalPacket    uint16
	currentPacket  uint16
	payloadLen     int
	payload        *Payload
}

func (p Packet) Marshal() []byte {
	message := make([]byte, 0x14) // build head

	// Head Flag, VERSION, RESERVED01, RESERVED02
	copy(message, regularPacketHeader)

	// SESSION ID
	binary.LittleEndian.PutUint32(message[0x04:], p.sessionID)

	// SEQUENCE NUMBER
	binary.LittleEndian.PutUint32(message[0x08:], p.sequenceNumber)

	// Total Packet
	// if p.TotalPacket == 0 {
	//	p.TotalPacket = 1
	// }

	message[0x0c] = byte(p.totalPacket)

	// CurPacket
	// if p.CurrentPacket == 0 {
	//	p.CurrentPacket = 1
	// }

	message[0x0d] = byte(p.currentPacket)

	// Message Id
	binary.LittleEndian.PutUint16(message[0x0e:], p.messageID)

	// Data Length
	var payloadLen int

	if p.payload != nil {
		payloadLen = p.payload.Len()
	}

	payloadLen += len(payloadEOF)
	binary.LittleEndian.PutUint32(message[0x10:], uint32(payloadLen))

	if p.payload != nil {
		// DATA
		message = append(message, p.payload.Bytes()...)
		message = append(message, payloadEOF...)
	}

	return message
}

func (p Packet) LoadPayload(payload interface{}) (err error) {
	switch pl := payload.(type) {
	case nil:
		p.payload.Reset()
		p.payloadLen = 0

	case string:
		p.payload.Reset()
		p.payloadLen, err = p.payload.Write([]byte(pl))

	case []byte:
		p.payload.Reset()
		p.payloadLen, err = p.payload.Write(pl)

	case io.Reader:
		var l int64

		p.payload.Reset()
		l, err = io.Copy(p.payload, pl)
		p.payloadLen = int(l)

	default:
		var encode []byte

		encode, err = json.Marshal(payload)
		if err == nil {
			p.payload.Reset()
			p.payloadLen, err = p.payload.Write(encode)
		}
	}

	return err
}

func newPacket() *Packet {
	return &Packet{
		payload: NewPayload(),
	}
}

func packetUnmarshal(message []byte) (*Packet, error) {
	if !bytes.Equal(regularPacketHeader, message[:0x04]) {
		return nil, errors.New("invalid regular packet header")
	}

	packet := newPacket()

	packet.sessionID = binary.LittleEndian.Uint32(message[0x04:0x08])
	packet.sequenceNumber = binary.LittleEndian.Uint32(message[0x08:0x0c])
	packet.totalPacket = uint16(message[0x0c])
	packet.currentPacket = uint16(message[0x0d])
	packet.messageID = binary.LittleEndian.Uint16(message[0x0e:0x10])
	packet.payloadLen = int(binary.LittleEndian.Uint32(message[0x10:0x14]))
	packet.payload.Write(message[0x14:])

	return packet, nil
}

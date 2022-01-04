package v1

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/kihamo/boggart/protocols/serial"
)

type Packet struct {
	payload []byte
	address uint32
	crc     uint16
	command uint8

	lock sync.RWMutex
}

func NewPacket() *Packet {
	return &Packet{}
}

func (r *Packet) Address() uint32 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.address
}

func (r *Packet) WithAddress(address uint32) *Packet {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.address = address

	return r
}

func (r *Packet) Command() uint8 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.command
}

func (r *Packet) WithCommand(command uint8) *Packet {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.command = command

	return r
}

func (r *Packet) Payload() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.payload...)
}

func (r *Packet) PayloadAsBuffer() *Buffer {
	return NewBuffer(r.Payload())
}

func (r *Packet) WithPayload(payload []byte) *Packet {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.payload = payload

	return r
}

func (r *Packet) CRC() uint16 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.crc
}

func (r *Packet) MarshalBinary() (_ []byte, err error) {
	var buf []byte

	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, r.Address())

	packet := append(buf, r.Command())
	packet = append(packet, r.Payload()...)

	r.lock.Lock()
	r.crc = binary.LittleEndian.Uint16(serial.GenerateCRC16(packet))
	r.lock.Unlock()

	buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, r.CRC())
	packet = append(packet, buf...)

	return packet, nil
}

func (r *Packet) UnmarshalBinary(data []byte) (err error) {
	l := len(data)

	if l < 7 {
		return errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	r.address = binary.LittleEndian.Uint32(data[:4])
	r.command = data[4]
	r.payload = data[5 : l-2]
	r.crc = binary.LittleEndian.Uint16(data[l-2:])

	// check crc
	crc := serial.GenerateCRC16(data[:l-2])
	if !bytes.Equal(crc, data[l-2:]) {
		return errors.New("wrong CRC16 of response packet " +
			hex.EncodeToString(data) + " have " +
			hex.EncodeToString(data[l-2:]) + " want " +
			hex.EncodeToString(crc))
	}

	return nil
}

func (r *Packet) String() string {
	data, _ := r.MarshalBinary()
	return hex.EncodeToString(data)
}

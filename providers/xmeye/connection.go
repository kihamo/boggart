package xmeye

import (
	"sync"
	"sync/atomic"

	protocol "github.com/kihamo/boggart/protocols/connection"
)

type connection struct {
	protocol.Conn
	sync.Mutex

	_              [4]byte
	sessionID      uint32
	sequenceNumber uint32
}

func (c *connection) SessionID() uint32 {
	return atomic.LoadUint32(&c.sessionID)
}

func (c *connection) SessionIDAsString() string {
	session := Uint32(c.SessionID())
	b, _ := session.MarshalJSON()

	return string(b)
}

func (c *connection) SetSessionID(id uint32) {
	atomic.StoreUint32(&c.sessionID, id)
}

func (c *connection) SequenceNumber() uint32 {
	return atomic.LoadUint32(&c.sequenceNumber)
}

func (c *connection) InSequenceNumber() {
	atomic.AddUint32(&c.sequenceNumber, 1)
}

func (c *connection) send(packet *packet) error {
	c.Lock()
	defer c.Unlock()

	packet.SessionID = c.SessionID()
	packet.SequenceNumber = c.SequenceNumber()

	if _, err := c.Write(packet.Marshal()); err != nil {
		return err
	}

	c.InSequenceNumber()

	return nil
}

func (c *connection) receive() (*packet, error) {
	c.Lock()
	defer c.Unlock()

	head := make([]byte, 0x14) // read head
	if _, err := c.Conn.Read(head); err != nil {
		return nil, err
	}

	// save session id
	packet, err := packetUnmarshal(head)
	if err != nil {
		return nil, err
	}

	if packet.PayloadLen == 0 {
		return packet, nil
	}

	bufSize := defaultPayloadBuffer
	if bufSize > packet.PayloadLen {
		bufSize = packet.PayloadLen
	}
	buf := make([]byte, bufSize)

	for {
		n, err := c.Read(buf)
		if err != nil {
			return packet, err
		}

		n, err = packet.Payload.Write(buf[:n])
		if err != nil {
			return packet, err
		}

		// чтобы не прочитать следующий пакет урезаем буфер (актуально для режима потока, там пакеты идут один за другим)
		if delta := packet.PayloadLen - packet.Payload.Len(); delta < bufSize {
			buf = make([]byte, delta)
		}

		if packet.Payload.Len() >= packet.PayloadLen {
			break
		}
	}

	return packet, nil
}

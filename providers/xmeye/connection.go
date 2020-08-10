package xmeye

import (
	"sync"
	"sync/atomic"

	protocol "github.com/kihamo/boggart/protocols/connection"
)

type connection struct {
	_ int64

	sessionID      uint32
	sequenceNumber uint32

	protocol.Connection

	lock sync.Mutex
}

func (c *connection) SessionID() uint32 {
	return atomic.LoadUint32(&c.sessionID)
}

func (c *connection) SessionIDAsString() (id string) {
	session := Uint32(c.SessionID())

	if b, err := session.MarshalJSON(); err == nil {
		id = string(b)
	}

	return id
}

func (c *connection) SetSessionID(id uint32) {
	atomic.StoreUint32(&c.sessionID, id)
}

func (c *connection) SequenceNumber() uint32 {
	return atomic.LoadUint32(&c.sequenceNumber)
}

func (c *connection) IncrementSequenceNumber() {
	atomic.AddUint32(&c.sequenceNumber, 1)
}

func (c *connection) send(packet *Packet) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	packet.sessionID = c.SessionID()
	packet.sequenceNumber = c.SequenceNumber()

	if _, err := c.Write(packet.Marshal()); err != nil {
		return err
	}

	c.IncrementSequenceNumber()

	return nil
}

func (c *connection) receive() (*Packet, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	head := make([]byte, 0x14) // read head
	if _, err := c.Connection.Read(head); err != nil {
		return nil, err
	}

	// save session id
	packet, err := packetUnmarshal(head)
	if err != nil {
		return nil, err
	}

	if packet.payloadLen == 0 {
		return packet, nil
	}

	bufSize := defaultPayloadBuffer

	if bufSize > packet.payloadLen {
		bufSize = packet.payloadLen
	}

	buf := make([]byte, bufSize)

	for {
		n, err := c.Read(buf)
		if err != nil {
			return packet, err
		}

		_, err = packet.payload.Write(buf[:n])
		if err != nil {
			return packet, err
		}

		// чтобы не прочитать следующий пакет урезаем буфер (актуально для режима потока, там пакеты идут один за другим)
		if delta := packet.payloadLen - packet.payload.Len(); delta < bufSize {
			buf = make([]byte, delta)
		}

		if packet.payload.Len() >= packet.payloadLen {
			break
		}
	}

	return packet, nil
}

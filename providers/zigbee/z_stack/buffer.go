package z_stack

import (
	"bytes"
	"encoding/binary"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(data),
	}
}

func (b *Buffer) ReadUint8() uint8 {
	value, _ := b.ReadByte()
	return value
}

func (b *Buffer) ReadUint16() uint16 {
	return binary.LittleEndian.Uint16(b.Next(2))
}

func (b *Buffer) ReadUint32() uint32 {
	return binary.LittleEndian.Uint32(b.Next(4))
}

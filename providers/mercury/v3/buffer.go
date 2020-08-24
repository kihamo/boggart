package v3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

var (
	zero = []byte{255, 255, 255, 255}
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(data),
	}
}

func (b *Buffer) readUint32(data []byte) uint32 {
	if len(data) == 3 {
		data = []byte{data[0], 0, data[1], data[2]}
	}

	if bytes.Equal(data, zero) {
		return 0
	}

	return binary.LittleEndian.Uint32([]byte{data[2], data[3], data[0], data[1]})
}

func (b *Buffer) ReadUint8() uint8 {
	value, _ := b.ReadByte()
	return value
}

func (b *Buffer) ReadUint16() uint16 {
	return binary.LittleEndian.Uint16(b.Next(2))
}

func (b *Buffer) ReadUint32() uint32 {
	return b.readUint32(b.Next(4))
}

func (b *Buffer) ReadUint32By3Byte() uint32 {
	return b.readUint32(b.Next(3))
}

func (b *Buffer) ReadFirmwareVersion() string {
	data := b.Next(3)
	return fmt.Sprintf("%d.%d.%d", data[0], data[1], data[2])
}

func (b *Buffer) ReadSerialNumber() string {
	data := b.Next(4)
	return fmt.Sprintf("%d%d%d%d", data[0], data[1], data[2], data[3])
}

func (b *Buffer) ReadBuildDate() time.Time {
	data := b.Next(3)
	return time.Date(2000+int(data[2]), time.Month(int(data[1])), int(data[0]), 0, 0, 0, 0, time.UTC)
}

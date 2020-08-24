package v1

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(data),
	}
}

func (b *Buffer) ReadBCD(n int) uint64 {
	result, _ := strconv.ParseUint(hex.EncodeToString(b.Next(n)), 10, 64)
	return result
}

func (b *Buffer) ReadUint8() uint8 {
	result, _ := b.ReadByte()
	return result
}

func (b *Buffer) ReadUint16() uint16 {
	return binary.LittleEndian.Uint16(b.Next(2))
}

func (b *Buffer) ReadUint32() uint32 {
	return binary.LittleEndian.Uint32(b.Next(4))
}

func (b *Buffer) ReadBool() bool {
	return b.ReadUint8() != 0
}

/*
	Length: 6
	Format: BCD
	Value: hh-mm-ss-dd-mon-yy
*/
func (b *Buffer) ReadTimeDate(loc *time.Location) time.Time {
	hh, mm, ss := b.ReadBCD(1), b.ReadBCD(1), b.ReadBCD(1)
	dd, mon, yy := b.ReadBCD(1), b.ReadBCD(1), b.ReadBCD(1)

	if dd == 0 && mon == 0 && yy == 0 {
		return time.Time{}
	}

	return time.Date(2000+int(yy), time.Month(mon), int(dd), int(hh), int(mm), int(ss), 0, loc)
}

/*
	Length: 7
	Format: BCD
	Value: dow-hh-mm-ss-dd-mon-yy
*/
func (b *Buffer) ReadTimeDateWithDayOfWeek(loc *time.Location) time.Time {
	b.ReadByte() // skip day of week
	return b.ReadTimeDate(loc)
}

func (b *Buffer) ReadDate() time.Time {
	dd, mon, yy := b.ReadBCD(1), b.ReadBCD(1), b.ReadBCD(1)

	return time.Date(2000+int(yy), time.Month(mon), int(dd), 0, 0, 0, 0, time.UTC)
}

func (b *Buffer) ReadCount() uint64 {
	return b.ReadBCD(4) * 10
}

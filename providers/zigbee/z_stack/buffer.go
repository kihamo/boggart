package z_stack

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/kihamo/boggart/protocols/serial"
)

// https://github.com/zigbeer/zcl-packet/wiki/6.-Appendix

const (
	DataTypeNoData       uint8 = 0
	DataTypeData8        uint8 = 8
	DataTypeData16       uint8 = 9
	DataTypeData24       uint8 = 10
	DataTypeData32       uint8 = 11
	DataTypeData40       uint8 = 12
	DataTypeData48       uint8 = 13
	DataTypeData56       uint8 = 14
	DataTypeData64       uint8 = 15
	DataTypeBoolean      uint8 = 16
	DataTypeBitMap8      uint8 = 24
	DataTypeBitMap16     uint8 = 25
	DataTypeBitMap24     uint8 = 26
	DataTypeBitMap32     uint8 = 27
	DataTypeBitMap40     uint8 = 28
	DataTypeBitMap48     uint8 = 29
	DataTypebitmap56     uint8 = 30
	DataTypebitmap64     uint8 = 31
	DataTypeUint8        uint8 = 32
	DataTypeUint16       uint8 = 33
	DataTypeUint24       uint8 = 34
	DataTypeUint32       uint8 = 35
	DataTypeUint40       uint8 = 36
	DataTypeUint48       uint8 = 37
	DataTypeUint56       uint8 = 38
	DataTypeUint64       uint8 = 39
	DataTypeInt8         uint8 = 40
	DataTypeInt16        uint8 = 41
	DataTypeInt24        uint8 = 42
	DataTypeInt32        uint8 = 43
	DataTypeEnum8        uint8 = 48
	DataTypeEnum16       uint8 = 49
	DataTypeSinglePrec   uint8 = 57
	DataTypeDoublePrec   uint8 = 58
	DataTypeOctetStr     uint8 = 65
	DataTypeCharStr      uint8 = 66
	DataTypeLongOctetStr uint8 = 67
	DataTypeLongCharStr  uint8 = 68
	DataTypeArray        uint8 = 72
	DataTypeStruct       uint8 = 76
	DataTypeSet          uint8 = 80
	DataTypeBag          uint8 = 81
	DataTypeTod          uint8 = 224
	DataTypeDate         uint8 = 225
	DataTypeUtc          uint8 = 226
	DataTypeClusterId    uint8 = 232
	DataTypeAttrId       uint8 = 233
	DataTypeBacOid       uint8 = 234
	DataTypeIEEEAddr     uint8 = 240
	DataTypeSecKey       uint8 = 241
	DataTypeUnknown      uint8 = 255
)

type TypeStruct struct {
	Count    uint16              // numElms
	Elements []TypeStructElement // structElms
}

type TypeStructElement struct {
	Type  uint8
	Value interface{}
}

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(data),
	}
}

func (b *Buffer) Frame(cmd0, cmd1 uint16) *Frame {
	f := NewFrame(cmd0, cmd1)
	f.SetDataAsBuffer(b)

	return f
}

func (b *Buffer) ReadByType(t uint8) interface{} {
	switch t {
	case DataTypeBoolean:
		return b.ReadBoolean()
	case DataTypeUint8:
		return b.ReadUint8()
	case DataTypeUint16:
		return b.ReadUint16()
	case DataTypeUint40:
		lsb, msb := b.ReadUint40()
		return []interface{}{lsb, msb}
	case DataTypeCharStr:
		return b.ReadCharStr()
	case DataTypeStruct:
		return b.ReadStruct()
	case DataTypeIEEEAddr:
		return b.ReadIEEEAddr()
	default:
		panic("unknown type of buffer " + fmt.Sprintf("%d %s", t, b.Bytes()))
	}

	return nil
}

func (b *Buffer) ReadBoolean() bool {
	return b.ReadUint8() != 0
}

func (b *Buffer) WriteBoolean(value bool) {
	if value {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func (b *Buffer) ReadUint8() uint8 {
	value, _ := b.ReadByte()
	return value
}

func (b *Buffer) WriteUint8(value uint8) {
	b.WriteByte(value)
}

func (b *Buffer) ReadUint16() uint16 {
	return binary.LittleEndian.Uint16(b.Next(2))
}

func (b *Buffer) WriteUint16(value uint16) {
	convert := make([]byte, 2)
	binary.LittleEndian.PutUint16(convert, value)

	b.Write(convert)
}

func (b *Buffer) ReadUint32() uint32 {
	return binary.LittleEndian.Uint32(b.Next(4))
}

func (b *Buffer) WriteUint32(value uint32) {
	convert := make([]byte, 4)
	binary.LittleEndian.PutUint32(convert, value)

	b.Write(convert)
}

func (b *Buffer) ReadUint40() (uint32, uint8) {
	return b.ReadUint32(), b.ReadUint8()
}

func (b *Buffer) ReadCharStr() string {
	l := int(b.ReadUint8())
	return string(b.Next(l))
}

func (b *Buffer) ReadStruct() TypeStruct {
	s := TypeStruct{
		Count: b.ReadUint16(),
	}

	s.Elements = make([]TypeStructElement, 0, s.Count)

	for i := uint16(0); i < s.Count; i++ {
		t := b.ReadUint8()

		s.Elements = append(s.Elements, TypeStructElement{
			Type:  t,
			Value: b.ReadByType(t),
		})
	}

	return s
}

func (b *Buffer) ReadIEEEAddr() []byte {
	return serial.Reverse(b.Next(8))
}

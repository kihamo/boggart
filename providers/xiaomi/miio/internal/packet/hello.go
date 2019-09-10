package packet

import (
	"bytes"
)

type Hello struct {
	Base
}

func NewHello() *Hello {
	packet := &Hello{}

	packet.Header.Length = 0x0020
	packet.Header.Unknown1 = bytes.Repeat([]byte{0xff}, 4)
	packet.Header.DeviceID = bytes.Repeat([]byte{0xff}, 4)
	packet.Header.Stamp = 0xffffffff
	packet.Header.Checksum = bytes.Repeat([]byte{0xff}, 16)

	return packet
}

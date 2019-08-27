package v3

import (
	"encoding/hex"

	"github.com/kihamo/boggart/components/boggart/protocols/serial"
)

type requestCode byte
type accessLevel byte

const (
	RequestCodeChannelTest     requestCode = 0x0
	RequestCodeChannelOpen     requestCode = 0x1
	RequestCodeChannelClose    requestCode = 0x2
	RequestCodeWriteParameters requestCode = 0x3
	RequestCodeReadData        requestCode = 0x4
	RequestCodeReadArray       requestCode = 0x5
	RequestCodeReadByAddress   requestCode = 0x6
	RequestCodeWriteByAddress  requestCode = 0x7
	RequestCodeReadParameter   requestCode = 0x8

	AccessLevel1 accessLevel = 0x1
	AccessLevel2 accessLevel = 0x2
)

type Request struct {
	Address            byte
	Code               requestCode
	ParameterCode      *byte
	ParameterExtension *byte
	Parameters         []byte
	CRC                []byte
}

func (r *Request) Bytes() []byte {
	packet := append([]byte{r.Address}, byte(r.Code))

	if r.ParameterCode != nil {
		packet = append(packet, *r.ParameterCode)
	}

	if r.ParameterExtension != nil {
		packet = append(packet, *r.ParameterExtension)
	}

	if len(r.Parameters) > 0 {
		packet = append(packet, r.Parameters...)
	}

	if len(r.CRC) == 0 {
		r.CRC = serial.GenerateCRC16(packet)
	}

	packet = append(packet, r.CRC...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

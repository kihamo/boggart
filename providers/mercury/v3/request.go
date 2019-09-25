package v3

import (
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

type requestCode int
type accessLevel int

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

	ParamCodeSerialNumberAndBuildDate       byte = 0x0
	ParamCodeForceReadParameters            byte = 0x01
	ParamCodeTransformationCoefficient      byte = 0x02
	ParamCodeVersion                        byte = 0x03
	ParamCodeAddress                        byte = 0x05
	ParamCodeAuxiliaryParameters3           byte = 0x11
	ParamCodeType                           byte = 0x12
	ParamCodeSelfDiagnostics                byte = 0x0A
	ParamCodeLocation                       byte = 0x0B
	ParamCodeAuxiliaryParameters12          byte = 0x16
	ParamCodeTarifficatorStatus             byte = 0x17
	ParamCodeLoadManager                    byte = 0x18
	ParamCodeLimitPower                     byte = 0x19
	ParamCodeMultiplierTimeoutMainInterface byte = 0x1D
	ParamCodeCRC16                          byte = 0x26
)

type Request struct {
	Address            byte
	Code               requestCode
	ParameterCode      *byte
	ParameterExtension *byte
	Parameters         []byte
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

	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

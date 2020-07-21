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
	address            byte
	code               requestCode
	parameterCode      *byte
	parameterExtension *byte
	parameters         []byte
}

func NewRequest(address byte, code requestCode) *Request {
	return &Request{
		address: address,
		code:    code,
	}
}

func (r *Request) WithParameterCode(code uint8) *Request {
	r.parameterCode = &[]byte{byte(code)}[0]
	return r
}

func (r *Request) WithParameterExtension(ext uint8) *Request {
	r.parameterExtension = &[]byte{byte(ext)}[0]
	return r
}

func (r *Request) WithParameters(params []byte) *Request {
	r.parameters = params
	return r
}

func (r *Request) Bytes() []byte {
	packet := append([]byte{r.address}, byte(r.code))

	if r.parameterCode != nil {
		packet = append(packet, *r.parameterCode)
	}

	if r.parameterExtension != nil {
		packet = append(packet, *r.parameterExtension)
	}

	if len(r.parameters) > 0 {
		packet = append(packet, r.parameters...)
	}

	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

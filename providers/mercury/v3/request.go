package v3

import (
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

const (
	RequestCodeChannelTest     uint8 = 0x0
	RequestCodeChannelOpen     uint8 = 0x1
	RequestCodeChannelClose    uint8 = 0x2
	RequestCodeWriteParameters uint8 = 0x3
	RequestCodeReadData        uint8 = 0x4
	RequestCodeReadArray       uint8 = 0x5
	RequestCodeReadByAddress   uint8 = 0x6
	RequestCodeWriteByAddress  uint8 = 0x7
	RequestCodeReadParameter   uint8 = 0x8

	AccessLevel1 uint8 = 0x1
	AccessLevel2 uint8 = 0x2

	ParamCodeSerialNumberAndBuildDate       uint8 = 0x00
	ParamCodeForceReadParameters            uint8 = 0x01
	ParamCodeTransformationCoefficient      uint8 = 0x02
	ParamCodeVersion                        uint8 = 0x03
	ParamCodeAddress                        uint8 = 0x05
	ParamCodeSelfDiagnostics                uint8 = 0x0A
	ParamCodeLocation                       uint8 = 0x0B
	ParamCodeAuxiliaryParameters3           uint8 = 0x11
	ParamCodeType                           uint8 = 0x12
	ParamCodeAuxiliaryParameters12          uint8 = 0x16
	ParamCodeTarifficatorStatus             uint8 = 0x17
	ParamCodeLoadManager                    uint8 = 0x18
	ParamCodeLimitPower                     uint8 = 0x19
	ParamCodeMultiplierTimeoutMainInterface uint8 = 0x1D
	ParamCodeCRC16                          uint8 = 0x26
)

type Request struct {
	address            uint8
	code               uint8
	parameterCode      *uint8
	parameterExtension *uint8
	parameters         []byte
}

func NewRequest(code uint8) *Request {
	return &Request{
		code: code,
	}
}

func (r *Request) WithAddress(address uint8) *Request {
	r.address = address
	return r
}

func (r *Request) WithParameterCode(code uint8) *Request {
	r.parameterCode = &[]uint8{code}[0]
	return r
}

func (r *Request) WithParameterExtension(ext uint8) *Request {
	r.parameterExtension = &[]uint8{ext}[0]
	return r
}

func (r *Request) WithParameters(params []byte) *Request {
	r.parameters = params
	return r
}

func (r *Request) Bytes() []byte {
	packet := append([]byte{r.address}, r.code)

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

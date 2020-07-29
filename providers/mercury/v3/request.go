package v3

import (
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

const (
	CodeChannelTest     = 0x0
	CodeChannelOpen     = 0x1
	CodeChannelClose    = 0x2
	CodeWriteParameters = 0x3
	CodeReadData        = 0x4
	CodeReadArray       = 0x5
	CodeReadByAddress   = 0x6
	CodeWriteByAddress  = 0x7
	CodeReadParameter   = 0x8

	AccessLevel1 = 0x1
	AccessLevel2 = 0x2

	ParamCodeSerialNumberAndBuildDate       = 0x00
	ParamCodeForceReadParameters            = 0x01
	ParamCodeTransformationCoefficient      = 0x02
	ParamCodeVersion                        = 0x03
	ParamCodeAddress                        = 0x05
	ParamCodeSelfDiagnostics                = 0x0A
	ParamCodeLocation                       = 0x0B
	ParamCodeAuxiliaryParameters3           = 0x11
	ParamCodeType                           = 0x12
	ParamCodeAuxiliaryParameters12          = 0x16
	ParamCodeTarifficatorStatus             = 0x17
	ParamCodeLoadManager                    = 0x18
	ParamCodeLimitPower                     = 0x19
	ParamCodeMultiplierTimeoutMainInterface = 0x1D
	ParamCodeCRC16                          = 0x26
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

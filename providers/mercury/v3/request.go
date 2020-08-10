package v3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"

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
	crc                uint16
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) Address() uint8 {
	return r.address
}

func (r *Request) WithAddress(address uint8) *Request {
	r.address = address
	return r
}

func (r *Request) Code() uint8 {
	return r.code
}

func (r *Request) WithCode(code uint8) *Request {
	r.code = code
	return r
}

func (r *Request) ParameterCode() *uint8 {
	return r.parameterCode
}

func (r *Request) WithParameterCode(code uint8) *Request {
	r.parameterCode = &[]uint8{code}[0]
	return r
}

func (r *Request) ParameterExtension() *uint8 {
	return r.parameterExtension
}

func (r *Request) WithParameterExtension(ext uint8) *Request {
	r.parameterExtension = &[]uint8{ext}[0]
	return r
}

func (r *Request) Parameters() []byte {
	return append([]byte(nil), r.parameters...)
}

func (r *Request) CRC() uint16 {
	return r.crc
}

func (r *Request) WithParameters(params []byte) *Request {
	r.parameters = params
	return r
}

func (r *Request) MarshalBinary() (_ []byte, err error) {
	packet := append([]byte{r.Address()}, r.Code())

	if code := r.ParameterCode(); code != nil {
		packet = append(packet, *code)
	}

	if ext := r.ParameterExtension(); ext != nil {
		packet = append(packet, *ext)
	}

	if params := r.Parameters(); len(params) > 0 {
		packet = append(packet, params...)
	}

	crc := serial.GenerateCRC16(packet)
	r.crc = binary.LittleEndian.Uint16(crc)
	packet = append(packet, crc...)

	return packet, nil
}

func (r *Request) UnmarshalBinary(data []byte) (err error) {
	l := len(data)

	if len(data) < 4 {
		return errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	r.address = data[0]
	r.code = data[1]

	if len(data) > 4 {
		r.parameterCode = &data[2]
	}

	if len(data) > 5 {
		r.parameterExtension = &data[3]
	}

	if len(data) > 6 {
		r.parameters = data[4 : l-2]
	}

	r.crc = binary.LittleEndian.Uint16(data[l-2:])

	// check crc
	crc := serial.GenerateCRC16(data[:l-2])
	if !bytes.Equal(crc, data[l-2:]) {
		return errors.New("error CRC16 of response packet " +
			hex.EncodeToString(data) + " have " +
			hex.EncodeToString(data[l-2:]) + " want " +
			hex.EncodeToString(crc))
	}

	return nil
}

func (r *Request) String() string {
	data, _ := r.MarshalBinary()
	return hex.EncodeToString(data)
}

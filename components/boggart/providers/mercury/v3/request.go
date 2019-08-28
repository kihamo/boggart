package v3

import (
	"encoding/hex"

	"github.com/kihamo/boggart/components/boggart/protocols/serial"
)

type requestCode int
type accessLevel int
type array int
type month int
type tariff int

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
	ParamCodeType                           byte = 0x12
	ParamCodeSelfDiagnostics                byte = 0x0A
	ParamCodeLocation                       byte = 0x0B
	ParamCodeTarifficatorStatus             byte = 0x17
	ParamCodeLoadManager                    byte = 0x18
	ParamCodeLimitPower                     byte = 0x19
	ParamCodeMultiplierTimeoutMainInterface byte = 0x1D
	ParamCodeCRC16                          byte = 0x26

	ArrayReset                 array = 0x00
	ArrayCurrentYear           array = 0x10
	ArrayPreviousYear          array = 0x20
	ArrayMonth                 array = 0x30
	ArrayCurrentDay            array = 0x40
	ArrayPreviousDay           array = 0x50
	ArrayActiveEnergy          array = 0x60
	ArrayBeginningCurrentYear  array = 0x90
	ArrayBeginningPreviousYear array = 0xA0
	ArrayBeginningMonth        array = 0xB0
	ArrayBeginningCurrentDay   array = 0xC0
	ArrayBeginningPreviousDay  array = 0xD0

	MonthJanuary   month = 0x01
	MonthFebruary  month = 0x02
	MonthMarch     month = 0x03
	MonthApril     month = 0x04
	MonthMay       month = 0x05
	MonthJune      month = 0x06
	MonthJuly      month = 0x07
	MonthAugust    month = 0x08
	MonthSeptember month = 0x09
	MonthOctober   month = 0x0A
	MonthNovember  month = 0x0B
	MonthDecember  month = 0x0C

	TariffAll tariff = 0x0
	Tariff1   tariff = 0x1
	Tariff2   tariff = 0x2
	Tariff3   tariff = 0x3
	Tariff4   tariff = 0x4
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

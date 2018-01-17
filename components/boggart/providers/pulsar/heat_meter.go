package pulsar

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
)

const (
	Input1 = iota + int64(1)
	Input2

	FunctionBadCommand    = 0x00
	FunctionReadMetrics   = 0x01
	FunctionReadTime      = 0x04
	FunctionWriteTime     = 0x05
	FunctionReadArchive   = 0x06
	FunctionReadSettings  = 0x0A
	FunctionWriteSettings = 0x0B

	ParamDaylightSavingTime = 0x0001 // uint16  | RW |    | признак автоперехода на летнее время 0 - выкл, 1 - вкл
	ParamPulseDuration      = 0x0003 // float32 | RW | мс | длительность импульса
	ParamPauseDuration      = 0x0004 // float32 | RW | мс | длительность паузы
	ParamVersion            = 0x0005 // uint16  | R  |    | версия прошивки
	ParamDiagnostics        = 0x0006
	ParamOperatingTime      = 0x000C // uint32  | RW | ч  | время наработки

	Channel3  = 0x00000004 // float32 | °C     | температура подачи
	Channel4  = 0x00000008 // float32 | °C     | температура обратки
	Channel5  = 0x00000010 // float32 | °C     | перепад температур
	Channel6  = 0x00000020 // float32 | гкал/ч | мощность
	Channel7  = 0x00000040 // float32 | гкал   | энергия
	Channel8  = 0x00000080 // float32 | м^3    | объем
	Channel9  = 0x00000100 // float32 | м^3/ч  | расход
	Channel10 = 0x00000200 // float32 | м^3    | импульсный вход 1
	Channel11 = 0x00000400 // float32 | м^3    | импульсный вход 2
	Channel12 = 0x00000800 // float32 | м^3    | импульсный вход 3
	Channel13 = 0x00001000 // float32 | м^3    | импульсный вход 4
	Channel14 = 0x00002000 // float32 | м3/ч   | расход по энергии === Channel6
	Channel20 = 0x00080000 // uint32  | ч      | время нормальной работы
)

type HeatMeter struct {
	address    []byte
	connection *rs485.Connection
}

func NewHeatMeter(address []byte, connection *rs485.Connection) *HeatMeter {
	return &HeatMeter{
		address:    address,
		connection: connection,
	}
}

func (d *HeatMeter) Address() []byte {
	return d.address
}

func (d *HeatMeter) Connection() *rs485.Connection {
	return d.connection
}

func (d *HeatMeter) Request(address []byte, function byte, data []byte) ([]byte, error) {
	var request []byte

	// device address
	request = append(request, address...)

	// function
	request = append(request, function)

	// length of packet
	l := len(request) + 1 + len(data) + 2 + 2
	request = append(request, byte(l))

	// data in
	request = append(request, data...)

	// request id
	requestId := rs485.GenerateRequestId()
	request = append(request, requestId...)

	// check sum CRC16
	request = append(request, rs485.GenerateCRC16(request)...)

	// fmt.Println("Request: ", request, p.ToString(request))

	response, err := d.connection.Request(request)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Response: ", response, p.ToString(response))

	l = len(response)
	if l < 10 {
		return nil, errors.New("Error length of response packet")
	}

	// check crc16
	crc16 := rs485.GenerateCRC16(response[:l-2])
	if bytes.Compare(response[l-2:], crc16) != 0 {
		return nil, errors.New("Error CRC16 of response packet")
	}

	// check id
	if bytes.Compare(response[l-(2+len(requestId)):l-2], requestId) != 0 {
		return nil, errors.New("Error ID of response packet")
	}

	// check error
	if response[4] == FunctionBadCommand {
		return nil, fmt.Errorf("HeatMeter returns error code #%d", response[6])
	}

	return response[6 : l-4], nil
}

func (d *HeatMeter) ReadMetrics(channel int64) ([][]byte, error) {
	bs := rs485.Pad(rs485.Reverse(big.NewInt(channel).Bytes()), 4)
	response, err := d.Request(d.address, FunctionReadMetrics, bs)
	if err != nil {
		return nil, err
	}

	metrics := math.Ceil(float64(len(response) / 4))
	result := make([][]byte, 0, int64(metrics))
	value := rs485.Reverse(response)
	for i := 0; i < len(value); i += 4 {
		result = append(result, value[i:i+4])
	}

	return result, nil
}

func (d *HeatMeter) ReadTime() (time.Time, error) {
	response, err := d.Request(d.address, FunctionReadTime, nil)
	if err != nil {
		return time.Now(), err
	}

	return time.Date(
		2000+int(response[0]),
		time.Month(response[1]),
		int(response[2]),
		int(response[3]),
		int(response[4]),
		int(response[5]),
		0,
		time.Now().Location()), nil
}

func (d *HeatMeter) ReadSettings(param int64) ([]byte, error) {
	bs := rs485.Pad(big.NewInt(param).Bytes(), 2)
	response, err := d.Request(d.address, FunctionReadSettings, bs)
	if err != nil {
		return nil, err
	}

	return rs485.Reverse(response), nil
}

func (d *HeatMeter) readMetricFloat32(channel int64) (float32, error) {
	value, err := d.ReadMetrics(channel)
	if err != nil {
		return -1, err
	}

	return rs485.ToFloat32(value[0]), nil
}

func (d *HeatMeter) TemperatureIn() (float32, error) {
	return d.readMetricFloat32(Channel3)
}

func (d *HeatMeter) TemperatureOut() (float32, error) {
	return d.readMetricFloat32(Channel4)
}

func (d *HeatMeter) TemperatureDelta() (float32, error) {
	return d.readMetricFloat32(Channel5)
}

func (d *HeatMeter) Power() (float32, error) {
	return d.readMetricFloat32(Channel6)
}

func (d *HeatMeter) Energy() (float32, error) {
	return d.readMetricFloat32(Channel7)
}

func (d *HeatMeter) Capacity() (float32, error) {
	return d.readMetricFloat32(Channel8)
}

func (d *HeatMeter) Consumption() (float32, error) {
	return d.readMetricFloat32(Channel9)
}

func (d *HeatMeter) PulseInput1() (float32, error) {
	return d.readMetricFloat32(Channel10)
}

func (d *HeatMeter) PulseInput2() (float32, error) {
	return d.readMetricFloat32(Channel11)
}

func (d *HeatMeter) PulseInput3() (float32, error) {
	return d.readMetricFloat32(Channel12)
}

func (d *HeatMeter) PulseInput4() (float32, error) {
	return d.readMetricFloat32(Channel13)
}

func (d *HeatMeter) DaylightSavingTime() (bool, error) {
	value, err := d.ReadSettings(ParamDaylightSavingTime)
	if err != nil {
		return false, err
	}

	return rs485.ToUint64(value) == 1, nil
}

func (d *HeatMeter) Diagnostics() ([]byte, error) {
	value, err := d.ReadSettings(ParamDiagnostics)
	if err != nil {
		return nil, err
	}

	// TODO: split result
	return value, nil
}

func (d *HeatMeter) Version() (uint16, error) {
	value, err := d.ReadSettings(ParamVersion)
	if err != nil {
		return 0, err
	}

	return uint16(rs485.ToUint64(value)), nil
}

func (d *HeatMeter) OperatingTime() (time.Duration, error) {
	value, err := d.ReadSettings(ParamOperatingTime)
	if err != nil {
		return -1, err
	}

	return time.Hour * time.Duration(rs485.ToUint64(value)), nil
}

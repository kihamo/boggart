package pulsar

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
)

const (
	FunctionBadCommand    = 0x00
	FunctionReadMetrics   = 0x01
	FunctionReadDatetime  = 0x04
	FunctionWriteTime     = 0x05
	FunctionReadArchive   = 0x06
	FunctionReadSettings  = 0x0A
	FunctionWriteSettings = 0x0B

	ArchiveTypeHourly  ArchiveType = 0x0001
	ArchiveTypeDaily   ArchiveType = 0x0002
	ArchiveTypeMonthly ArchiveType = 0x0003

	SettingsParamDaylightSavingTime  SettingsParam = 0x0001 // uint16  | RW |          | признак автоперехода на летнее время 0 - выкл, 1 - вкл
	SettingsParamPulseDuration       SettingsParam = 0x0003 // float32 | RW | мс       | длительность импульса
	SettingsParamPauseDuration       SettingsParam = 0x0004 // float32 | RW | мс       | длительность паузы
	SettingsParamVersion             SettingsParam = 0x0005 // uint16  | R  |          | версия прошивки
	SettingsParamDiagnostics         SettingsParam = 0x0006 // uint16  | R  |          | диагностика
	SettingsParamResetMCU            SettingsParam = 0x0007 // uint16  | R  |          | количество сбросов MCU
	SettingsParamBatteryVoltage      SettingsParam = 0x000A // float32 | R  | v        | напряжение батареи
	SettingsParamDeviceTemperature   SettingsParam = 0x000B // float32 | R  | c        | температура прибота
	SettingsParamOperatingTime       SettingsParam = 0x000C // uint32  | RW | ч        | время наработки
	SettingsParamErrorOperatingTime  SettingsParam = 0x000D // uint32  | RW | ч        | время наработки с ошибками
	SettingsParamPulse1Volume        SettingsParam = 0x0020 // float32 | RW | м3       | вес импульсного входа 1
	SettingsParamPulse1Duration      SettingsParam = 0x0021 // float32 | RW | мс       | длительность импульса импульсного входа 1
	SettingsParamPulse1PauseDuration SettingsParam = 0x0022 // float32 | RW | мс       | длительность паузы импульсного входа 1
	SettingsParamPulse2Volume        SettingsParam = 0x0023 // float32 | RW | м3       | вес импульсного входа 2
	SettingsParamPulse2Duration      SettingsParam = 0x0024 // float32 | RW | мс       | длительность импульса импульсного входа 2
	SettingsParamPulse2PauseDuration SettingsParam = 0x0025 // float32 | RW | мс       | длительность паузы импульсного входа 2
	SettingsParamPulse3Volume        SettingsParam = 0x0026 // float32 | RW | м3       | вес импульсного входа 3
	SettingsParamPulse3Duration      SettingsParam = 0x0027 // float32 | RW | мс       | длительность импульса импульсного входа 3
	SettingsParamPulse3PauseDuration SettingsParam = 0x0028 // float32 | RW | мс       | длительность паузы импульсного входа 3
	SettingsParamPulse4Volume        SettingsParam = 0x0029 // float32 | RW | м3       | вес импульсного входа 4
	SettingsParamPulse4Duration      SettingsParam = 0x002A // float32 | RW | мс       | длительность импульса импульсного входа 4
	SettingsParamPulse4PauseDuration SettingsParam = 0x002B // float32 | RW | мс       | длительность импульса импульсного входа 4
	SettingsParamOutputVolume        SettingsParam = 0x002C // float32 | RW | гкал/имп | вес импульса выхода
	SettingsParamOutputDuration      SettingsParam = 0x002D // float32 | RW | мс       | длительность импульса выхода

	Channel3  MetricsChannel = 0x00000004 // float32 | °C     | температура подачи
	Channel4  MetricsChannel = 0x00000008 // float32 | °C     | температура обратки
	Channel5  MetricsChannel = 0x00000010 // float32 | °C     | перепад температур
	Channel6  MetricsChannel = 0x00000020 // float32 | гкал/ч | мощность
	Channel7  MetricsChannel = 0x00000040 // float32 | гкал   | энергия
	Channel8  MetricsChannel = 0x00000080 // float32 | м^3    | объем
	Channel9  MetricsChannel = 0x00000100 // float32 | м^3/ч  | расход
	Channel10 MetricsChannel = 0x00000200 // float32 | м^3    | импульсный вход 1
	Channel11 MetricsChannel = 0x00000400 // float32 | м^3    | импульсный вход 2
	Channel12 MetricsChannel = 0x00000800 // float32 | м^3    | импульсный вход 3
	Channel13 MetricsChannel = 0x00001000 // float32 | м^3    | импульсный вход 4
	Channel14 MetricsChannel = 0x00002000 // float32 | м3/ч   | расход по энергии === Channel6
	Channel20 MetricsChannel = 0x00080000 // uint32  | ч      | время нормальной работы
)

type MetricsChannel int
type SettingsParam int
type ArchiveType int

type HeatMeter struct {
	address    []byte
	location   *time.Location
	connection *rs485.Connection
}

func (i MetricsChannel) toInt64() int64 {
	return int64(i)
}

func (i MetricsChannel) toBytes() []byte {
	return big.NewInt(i.toInt64()).Bytes()
}

func (i SettingsParam) toInt64() int64 {
	return int64(i)
}

func (i SettingsParam) toBytes() []byte {
	return big.NewInt(i.toInt64()).Bytes()
}

func (i ArchiveType) toInt64() int64 {
	return int64(i)
}

func (i ArchiveType) toBytes() []byte {
	return big.NewInt(i.toInt64()).Bytes()
}

func NewHeatMeter(address []byte, location *time.Location, connection *rs485.Connection) *HeatMeter {
	if location == nil {
		location = time.Now().Location()
	}

	return &HeatMeter{
		address:    address,
		location:   location,
		connection: connection,
	}
}

func (d *HeatMeter) Address() []byte {
	return d.address
}

func (d *HeatMeter) Connection() *rs485.Connection {
	return d.connection
}

func (d *HeatMeter) Request(function byte, data []byte) ([]byte, error) {
	var request []byte

	// device address
	request = append(request, d.address...)

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

	// fmt.Println("Request: ", request, hex.EncodeToString(request), " with function", strings.ToUpper(hex.EncodeToString([]byte{function})))

	response, err := d.connection.Request(request)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Response: ", response, hex.EncodeToString(response))

	l = len(response)
	if l < 10 {
		return nil, errors.New("error length of response packet")
	}

	// check crc16
	crc16 := rs485.GenerateCRC16(response[:l-2])
	if !bytes.Equal(response[l-2:], crc16) {
		return nil, errors.New("error CRC16 of response packet")
	}

	// check id
	if !bytes.Equal(response[l-(2+len(requestId)):l-2], requestId) {
		return nil, errors.New("error ID of response packet")
	}

	// check error
	if response[4] == FunctionBadCommand {
		return nil, fmt.Errorf("HeatMeter returns error code #%d", response[6])
	}

	return response[6 : l-4], nil
}

func (d *HeatMeter) readMetrics(channel MetricsChannel) (float32, error) {
	bs := rs485.Pad(rs485.Reverse(channel.toBytes()), 4)
	response, err := d.Request(FunctionReadMetrics, bs)
	if err != nil {
		return -1, err
	}

	return rs485.ToFloat32(rs485.Reverse(response)), nil
}

func (d *HeatMeter) Datetime() (time.Time, error) {
	response, err := d.Request(FunctionReadDatetime, nil)
	if err != nil {
		return time.Time{}, err
	}

	return BytesToTime(response, d.location), nil
}

/*
	Максимальная глубина архивов
	- Часовые 62 суток (1488 значений)
	- Суточные 6 месцев (184 суток)
	- Месячные 5 лет (60 значений)
*/
func (d *HeatMeter) readArchive(channel MetricsChannel, start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	/*
		DATE_START
		дата округляется прибором до ближайшей архивной записи слева, в некоторых ранних прошивках приборов
		нормировка архивов не производилась, поэтому желательно нормировку даты осуществлять софтом верхнего уровня
	*/
	switch t {
	case ArchiveTypeMonthly:
		start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, end.Location())
	case ArchiveTypeDaily:
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, end.Location())
	case ArchiveTypeHourly:
		start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, end.Location())
	}

	bs := rs485.Pad(rs485.Reverse(channel.toBytes()), 4)
	bs = append(bs, rs485.Pad(t.toBytes(), 2)...)
	bs = append(bs, TimeToBytes(start)...)
	bs = append(bs, TimeToBytes(end)...)

	response, err := d.Request(FunctionReadArchive, bs)
	if err != nil {
		return time.Time{}, nil, err
	}

	begin := BytesToTime(response[4:10], d.location)
	raw := rs485.Reverse(response[10:])
	values := make([]float32, 0)

	for i := 0; i < len(raw); i += 4 {
		values = append([]float32{rs485.ToFloat32(raw[i : i+4])}, values...)
	}

	return begin, values, nil
}

func (d *HeatMeter) readSettings(param SettingsParam) ([]byte, error) {
	bs := rs485.Pad(param.toBytes(), 2)
	response, err := d.Request(FunctionReadSettings, bs)
	if err != nil {
		return nil, err
	}

	return rs485.Reverse(response), nil
}

func (d *HeatMeter) TemperatureIn() (float32, error) {
	return d.readMetrics(Channel3)
}

func (d *HeatMeter) TemperatureInArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel3, start, end, t)
}

func (d *HeatMeter) TemperatureOut() (float32, error) {
	return d.readMetrics(Channel4)
}

func (d *HeatMeter) TemperatureOutArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel4, start, end, t)
}

func (d *HeatMeter) TemperatureDelta() (float32, error) {
	return d.readMetrics(Channel5)
}

func (d *HeatMeter) TemperatureDeltaArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel5, start, end, t)
}

func (d *HeatMeter) Power() (float32, error) {
	return d.readMetrics(Channel6)
}

func (d *HeatMeter) PowerArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel6, start, end, t)
}

func (d *HeatMeter) PowerByEnergy() (float32, error) {
	return d.readMetrics(Channel14)
}

func (d *HeatMeter) PowerByEnergyArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel14, start, end, t)
}

func (d *HeatMeter) Energy() (float32, error) {
	return d.readMetrics(Channel7)
}

func (d *HeatMeter) EnergyArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel7, start, end, t)
}

func (d *HeatMeter) Capacity() (float32, error) {
	return d.readMetrics(Channel8)
}

func (d *HeatMeter) CapacityArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel8, start, end, t)
}

func (d *HeatMeter) Consumption() (float32, error) {
	return d.readMetrics(Channel9)
}

func (d *HeatMeter) ConsumptionArchive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel9, start, end, t)
}

func (d *HeatMeter) PulseInput1() (float32, error) {
	return d.readMetrics(Channel10)
}

func (d *HeatMeter) PulseInput1Archive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel10, start, end, t)
}

func (d *HeatMeter) PulseInput2() (float32, error) {
	return d.readMetrics(Channel11)
}

func (d *HeatMeter) PulseInput2Archive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel11, start, end, t)
}

func (d *HeatMeter) PulseInput3() (float32, error) {
	return d.readMetrics(Channel12)
}

func (d *HeatMeter) PulseInput3Archive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel12, start, end, t)
}

func (d *HeatMeter) PulseInput4() (float32, error) {
	return d.readMetrics(Channel13)
}

func (d *HeatMeter) PulseInput4Archive(start, end time.Time, t ArchiveType) (time.Time, []float32, error) {
	return d.readArchive(Channel14, start, end, t)
}

func (d *HeatMeter) DaylightSavingTime() (bool, error) {
	value, err := d.readSettings(SettingsParamDaylightSavingTime)
	if err != nil {
		return false, err
	}

	return rs485.ToUint64(value) == 1, nil
}

func (d *HeatMeter) Diagnostics() ([]byte, error) {
	value, err := d.readSettings(SettingsParamDiagnostics)
	if err != nil {
		return nil, err
	}

	// TODO: split result
	return value, nil
}

func (d *HeatMeter) Version() (uint16, error) {
	value, err := d.readSettings(SettingsParamVersion)
	if err != nil {
		return 0, err
	}

	return uint16(rs485.ToUint64(value)), nil
}

func (d *HeatMeter) OperatingTime() (time.Duration, error) {
	value, err := d.readSettings(SettingsParamOperatingTime)
	if err != nil {
		return -1, err
	}

	return time.Hour * time.Duration(rs485.ToUint64(value)), nil
}

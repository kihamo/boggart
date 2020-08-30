package pulsar

import (
	"time"

	"github.com/kihamo/boggart/protocols/serial"
)

func (d *HeatMeter) readMetrics(channel MetricsChannel) (float32, error) {
	request := NewPacket().
		WithFunction(FunctionReadMetrics).
		WithPayload(serial.Pad(serial.Reverse(channel.toBytes()), 4))

	response, err := d.Invoke(request)

	if err != nil {
		return -1, err
	}

	return serial.ToFloat32(serial.Reverse(response.Payload())), nil
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

	bs := serial.Pad(serial.Reverse(channel.toBytes()), 4)
	bs = append(bs, serial.Pad(t.toBytes(), 2)...)
	bs = append(bs, TimeToBytes(start)...)
	bs = append(bs, TimeToBytes(end)...)

	request := NewPacket().
		WithFunction(FunctionReadArchive).
		WithPayload(bs)

	response, err := d.Invoke(request)
	if err != nil {
		return time.Time{}, nil, err
	}

	payload := response.Payload()

	begin := BytesToTime(payload[4:10], d.options.location)
	raw := serial.Reverse(payload[10:])
	values := make([]float32, 0)

	for i := 0; i < len(raw); i += 4 {
		values = append([]float32{serial.ToFloat32(raw[i : i+4])}, values...)
	}

	return begin, values, nil
}

func (d *HeatMeter) Datetime() (time.Time, error) {
	response, err := d.Invoke(NewPacket().WithFunction(FunctionReadDatetime))
	if err != nil {
		return time.Time{}, err
	}

	return BytesToTime(response.Payload(), d.options.location), nil
}

func (d *HeatMeter) ReadSettings(param SettingsParam) ([]byte, error) {
	request := NewPacket().
		WithFunction(FunctionReadSettings).
		WithPayload(serial.Pad(param.toBytes(), 2))

	response, err := d.Invoke(request)

	if err != nil {
		return nil, err
	}

	return serial.Reverse(response.Payload()), nil
}

func (d *HeatMeter) DaylightSavingTime() (bool, error) {
	value, err := d.ReadSettings(SettingsParamDaylightSavingTime)
	if err != nil {
		return false, err
	}

	return serial.ToUint64(value) == 1, nil
}

func (d *HeatMeter) Version() (uint16, error) {
	value, err := d.ReadSettings(SettingsParamVersion)
	if err != nil {
		return 0, err
	}

	return uint16(serial.ToUint64(value)), nil
}

func (d *HeatMeter) Diagnostics() ([]byte, error) {
	value, err := d.ReadSettings(SettingsParamDiagnostics)
	if err != nil {
		return nil, err
	}

	// TODO: split result
	return value, nil
}

func (d *HeatMeter) BatteryVoltage() (float32, error) {
	value, err := d.ReadSettings(SettingsParamBatteryVoltage)
	if err != nil {
		return -1, err
	}

	return serial.ToFloat32(serial.Reverse(value)), nil
}

func (d *HeatMeter) DeviceTemperature() ([]byte, error) {
	value, err := d.ReadSettings(SettingsParamDeviceTemperature)
	if err != nil {
		return nil, err
	}

	return value, err
}

func (d *HeatMeter) OperatingTime() (time.Duration, error) {
	value, err := d.ReadSettings(SettingsParamOperatingTime)
	if err != nil {
		return -1, err
	}

	return time.Hour * time.Duration(serial.ToUint64(value)), nil
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

package pulsar

import (
	"math/big"
	"strconv"
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

const (
	DiagnosticsEmptyBattery uint8 = 1 << iota
	DiagnosticsErrorReadWriteEEPROM
	DiagnosticsResetCounters
	DiagnosticsReserved1
	DiagnosticsReserved2
	DiagnosticsThermometerInBroke
	DiagnosticsThermometerOutBroke
	DiagnosticsNegativeTemperaturesDelta
)

type MetricsChannel int
type SettingsParam int
type ArchiveType int
type ErrorCode uint8

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

/*
ERROR_CODE - (uint8_t) код ошибки:
(0x01) - отсутствует запрашиваемый код функции;
(0x02) - ошибка в битовой маске запроса;
(0x03) - ошибочная длинна запроса;
(0x04) - отсутствует параметр
(0x05) - запись заблокирована, требуется авторизация;
(0x06) - записываемое значение (параметр) находится вне заданного
диапазона;
(0x07) - отсутствует запрашиваемый тип архива;
(0x08) – превышение максимального количества архивных значений за один
пакет;
*/
func (i ErrorCode) Error() string {
	switch i {
	case 0x01:
		return "function is empty"
	case 0x02:
		return "wrong bit mask"
	case 0x03:
		return "wrong length of request"
	case 0x04:
		return "parameter ie empty"
	case 0x05:
		return "write blocked, need authorization"
	case 0x06:
		return "range of value is wrong"
	case 0x07:
		return "type archive is empty"
	case 0x08:
		return "exceeding the maximum number of archived values"
	}

	return "unknown error code " + strconv.FormatUint(uint64(i), 10)
}

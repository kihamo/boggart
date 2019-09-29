package v1

import (
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

type requestCommand int

// https://github.com/mrkrasser/MercuryStats
// http://incotex-support.blogspot.ru/2016/05/blog-post.html

const (
	// группа сетевых команд установки
	RequestCommandWriteDisplayMode = 0x09 // установка режима индикации
	RequestCommandWriteDisplayTime = 0x0D // установка времени индикации

	// группа сетевых команд чтения
	RequestCommandReadAddressGroup         = 0x20 // чтение группового адреса счетчика
	RequestCommandReadDatetime             = 0x21 // чтение установленных даты и времени
	RequestCommandReadPowerMaximum         = 0x22 // чтение лимита мощности
	RequestCommandReadEnergyMaximum        = 0x23 // чтение лимита энергии
	RequestCommandReadDaylightSavingTime   = 0x24 // чтение значения флага сезонного времени
	RequestCommandReadTimeCorrection       = 0x25 // чтение величины коррекции времени
	RequestCommandReadPowerCurrent         = 0x26 // чтение текущей мощности в нагрузке
	RequestCommandReadPowerCounters        = 0x27 // чтение содержимого тарифных аккумуляторов
	RequestCommandReadVersion              = 0x28 // чтение версии ПО
	RequestCommandReadBatteryVoltage       = 0x29 // чтение напряжения встроенной батарейки
	RequestCommandReadDisplayMode          = 0x2A // чтение режима индикации
	RequestCommandReadLastPowerOffDatetime = 0x2B // чтение времени последнего отключения напряжения
	RequestCommandReadLastPowerOnDatetime  = 0x2C // чтение времени последнего включения напряжения
	RequestCommandReadSerialNumber         = 0x2F // чтение серийного номера
	RequestCommandReadTariffCount          = 0x2E // чтение количества действующих тарифов

	// группа дополнительных сетевых команд чтения
	RequestCommandReadHolidays         = 0x30 // чтение таблицы праздничных дней
	RequestCommandReadMonthlyStat      = 0x32 // чтение месячных срезов
	RequestCommandReadEventsPowerOnOff = 0x34 // чтение буфера событий включений/выключения
	RequestCommandReadCurrentTariff    = 0x60 // чтение тарифа
	RequestCommandReadLastCloseCap     = 0x62 // чтение времение последнего закрытия крышки счетчика
	RequestCommandReadParamsCurrent    = 0x63 // чтение текущих параметров сети U, I, P
	RequestCommandReadWordType         = 0x65 // чтение слова исполнения
	RequestCommandReadMakeDate         = 0x66 // чтение даты изготовления
	RequestCommandReadDisplayTime      = 0x67 // чтение времени индекации
	RequestCommandReadWorkingTime      = 0x69 // чтение времени наработки
)

type Request struct {
	Address []byte
	Command requestCommand
	Payload []byte
}

func (r *Request) Bytes() []byte {
	packet := append(r.Address, byte(r.Command))
	packet = append(packet, r.Payload...)
	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

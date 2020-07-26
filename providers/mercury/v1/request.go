package v1

import (
	"encoding/binary"
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
	RequestCommandReadMaximum          = 0x33 // чтение месячных срезов
	RequestCommandReadEventsPowerOnOff = 0x34 // чтение буфера событий включений/выключения
	RequestCommandReadEventsOpenClose  = 0x35 // чтение буфера событий включений/выключения

	RequestCommandReadCurrentTariff = 0x60 // чтение тарифа
	RequestCommandReadLastOpenCap   = 0x61 // Чтение времени последнего вскрытия крышки счётчика
	RequestCommandReadLastCloseCap  = 0x62 // чтение времение последнего закрытия крышки счетчика
	RequestCommandReadParamsCurrent = 0x63 // чтение текущих параметров сети U, I, P
	RequestCommandReadMakeDate      = 0x66 // чтение даты изготовления
	RequestCommandReadDisplayTime   = 0x67 // чтение времени индекации
	RequestCommandReadWorkingTime   = 0x69 // чтение времени наработки
)

type Request struct {
	address uint32
	command requestCommand
	payload []byte
}

func NewRequest(command requestCommand) *Request {
	return &Request{
		command: command,
	}
}

func (r *Request) WithAddress(address uint32) *Request {
	r.address = address
	return r
}

func (r *Request) WithPayload(payload []byte) *Request {
	r.payload = payload
	return r
}

func (r *Request) Bytes() []byte {
	packet := make([]byte, 4)
	binary.LittleEndian.PutUint32(packet, r.address)

	packet = append(packet, byte(r.command))
	packet = append(packet, r.payload...)
	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

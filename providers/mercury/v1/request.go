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
	// 0x00 // установка нового сетевого адреса счетчика
	// 0x01 // установка нового группового адреса счетчика
	// 0x02 // установка внутренних часов и календаря счетчика
	// 0x03 // установка лимита мощности
	// 0x04 // установка лимита энергии за месяц
	// 0x05 // установка флага сезонного времени
	// 0x06 // установка величины коррекции времени
	// 0x07 // установка функции выходного оптрона
	// 0x08 // установка скорости обмена
	RequestCommandWriteDisplayMode = 0x09 // установка режима индикации
	// 0x0A // установка числа действующих тарифов
	// 0x0B // установка тарифа
	// 0x0C // сброс защёлки "напряжение батареи"
	RequestCommandWriteDisplayTime = 0x0D // установка времени индикации
	// 0x0F // установка режима лимита мощности
	// 0x10 // установка таблицы праздничных дней
	// 0x11 // установка таблицы переключений тарифных зон
	// 0x12 // сброс максимумов
	// 0x13 // сброс максимумов. Под перемычкой
	// 0x14 // сброс наработки батареи.

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
	RequestCommandReadOptocouplerFunction  = 0x2D // Чтение функции выходного оптрона
	RequestCommandReadTariffCount          = 0x2E // чтение количества действующих тарифов
	RequestCommandReadSerialNumber         = 0x2F // чтение серийного номера

	// группа дополнительных сетевых команд чтения
	RequestCommandReadHolidays          = 0x30 // чтение таблицы праздничных дней
	RequestCommandReadTariffZoneChanged = 0x31 // чтение таблицы переключений тарифных зон
	RequestCommandReadMonthlyStat       = 0x32 // чтение месячных срезов
	RequestCommandReadMaximum           = 0x33 // чтение месячных срезов
	RequestCommandReadEventsPowerOnOff  = 0x34 // чтение буфера событий включений/выключения
	RequestCommandReadEventsOpenClose   = 0x35 // чтение буфера событий включений/выключения
	// 0x36 // чтение буфера событий параметризации
	// 0x37 // чтение получасовых мощностей или суточных срезов
	// 0x38 // чтение месячных срезов реактивной энергии
	// 0x39 // чтение буфера событий качества электричества
	RequestCommandReadEventsRelay = 0x3A // чтение буфера событий реле

	RequestCommandReadCurrentTariff = 0x60 // чтение тарифа
	RequestCommandReadLastOpenCap   = 0x61 // Чтение времени последнего вскрытия крышки счётчика
	RequestCommandReadLastCloseCap  = 0x62 // чтение времение последнего закрытия крышки счетчика
	RequestCommandReadParamsCurrent = 0x63 // чтение текущих параметров сети U, I, P
	// 0x64 // чтение коэффициента коррекции хода часов Введена для чтения коэффициента коррекции без перемычки.
	// 0x65 // чтение слова исполнения
	RequestCommandReadMakeDate    = 0x66 // чтение даты изготовления
	RequestCommandReadDisplayTime = 0x67 // чтение времени индекации
	// 0x68 // чтение режима лимита мощности
	RequestCommandReadWorkingTime = 0x69 // чтение времени наработки
	// 0x6A // чтение режима доп. индикации
	// 0x6B // чтение времени последней парам. счётчика
	// 0x6C // чтение номера модема и системы и уровня сигнала
	// 0x6D // чтение режима управления реле
	// 0x6E // чтение потарифных лимитов Энергии (остатки)
	// 0x6F // чтение флага разрешения индикации под батарейкой
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

package v1

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/kihamo/boggart/protocols/serial"
)

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
	RequestCommandWriteDisplayMode uint8 = 0x09 // установка режима индикации
	// 0x0A // установка числа действующих тарифов
	// 0x0B // установка тарифа
	// 0x0C // сброс защёлки "напряжение батареи"
	RequestCommandWriteDisplayTime uint8 = 0x0D // установка времени индикации
	// 0x0F // установка режима лимита мощности
	// 0x10 // установка таблицы праздничных дней
	// 0x11 // установка таблицы переключений тарифных зон
	// 0x12 // сброс максимумов
	// 0x13 // сброс максимумов. Под перемычкой
	// 0x14 // сброс наработки батареи.

	// группа сетевых команд чтения
	RequestCommandReadAddressGroup         uint8 = 0x20 // чтение группового адреса счетчика
	RequestCommandReadDatetime             uint8 = 0x21 // чтение установленных даты и времени
	RequestCommandReadPowerMaximum         uint8 = 0x22 // чтение лимита мощности
	RequestCommandReadEnergyMaximum        uint8 = 0x23 // чтение лимита энергии
	RequestCommandReadDaylightSavingTime   uint8 = 0x24 // чтение значения флага сезонного времени
	RequestCommandReadTimeCorrection       uint8 = 0x25 // чтение величины коррекции времени
	RequestCommandReadPowerCurrent         uint8 = 0x26 // чтение текущей мощности в нагрузке
	RequestCommandReadPowerCounters        uint8 = 0x27 // чтение содержимого тарифных аккумуляторов
	RequestCommandReadVersion              uint8 = 0x28 // чтение версии ПО
	RequestCommandReadBatteryVoltage       uint8 = 0x29 // чтение напряжения встроенной батарейки
	RequestCommandReadDisplayMode          uint8 = 0x2A // чтение режима индикации
	RequestCommandReadLastPowerOffDatetime uint8 = 0x2B // чтение времени последнего отключения напряжения
	RequestCommandReadLastPowerOnDatetime  uint8 = 0x2C // чтение времени последнего включения напряжения
	RequestCommandReadOptocouplerFunction  uint8 = 0x2D // Чтение функции выходного оптрона
	RequestCommandReadTariffCount          uint8 = 0x2E // чтение количества действующих тарифов
	RequestCommandReadSerialNumber         uint8 = 0x2F // чтение серийного номера

	// группа дополнительных сетевых команд чтения
	RequestCommandReadHolidays          uint8 = 0x30 // чтение таблицы праздничных дней
	RequestCommandReadTariffZoneChanged uint8 = 0x31 // чтение таблицы переключений тарифных зон
	RequestCommandReadMonthlyStat       uint8 = 0x32 // чтение месячных срезов
	RequestCommandReadMaximum           uint8 = 0x33 // чтение месячных срезов
	RequestCommandReadEventsPowerOnOff  uint8 = 0x34 // чтение буфера событий включений/выключения
	RequestCommandReadEventsOpenClose   uint8 = 0x35 // чтение буфера событий включений/выключения
	// 0x36 // чтение буфера событий параметризации
	// 0x37 // чтение получасовых мощностей или суточных срезов
	// 0x38 // чтение месячных срезов реактивной энергии
	// 0x39 // чтение буфера событий качества электричества
	RequestCommandReadEventsRelay uint8 = 0x3A // чтение буфера событий реле

	RequestCommandReadCurrentTariff uint8 = 0x60 // чтение тарифа
	RequestCommandReadLastOpenCap   uint8 = 0x61 // Чтение времени последнего вскрытия крышки счётчика
	RequestCommandReadLastCloseCap  uint8 = 0x62 // чтение времение последнего закрытия крышки счетчика
	RequestCommandReadParamsCurrent uint8 = 0x63 // чтение текущих параметров сети U, I, P
	// 0x64 // чтение коэффициента коррекции хода часов Введена для чтения коэффициента коррекции без перемычки.
	RequestCommandReadModel       uint8 = 0x65 // чтение слова исполнения
	RequestCommandReadMakeDate    uint8 = 0x66 // чтение даты изготовления
	RequestCommandReadDisplayTime uint8 = 0x67 // чтение времени индекации
	// 0x68 // чтение режима лимита мощности
	RequestCommandReadWorkingTime     uint8 = 0x69 // чтение времени наработки
	RequestCommandReadDisplayModeExt  uint8 = 0x6A // чтение режима доп. индикации
	RequestCommandReadParamLastChange uint8 = 0x6B // чтение времени последней парам. счётчика
	// 0x6C // чтение номера модема и системы и уровня сигнала
	RequestCommandReadRelayMode                   uint8 = 0x6D // чтение режима управления реле
	RequestCommandReadPowerLimits                 uint8 = 0x6E // чтение потарифных лимитов Энергии (остатки)
	RequestCommandReadAllowIndicationUnderBattery uint8 = 0x6F // чтение флага разрешения индикации под батарейкой
)

type Request struct {
	address uint32
	command uint8
	payload []byte
}

func NewRequest(command uint8) *Request {
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

	packet = append(packet, r.command)
	packet = append(packet, r.payload...)
	packet = append(packet, serial.GenerateCRC16(packet)...)

	return packet
}

func (r *Request) String() string {
	return hex.EncodeToString(r.Bytes())
}

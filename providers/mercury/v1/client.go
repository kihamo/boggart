package v1

import (
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/connection"
)

// https://github.com/mrkrasser/MercuryStats
// http://incotex-support.blogspot.ru/2016/05/blog-post.html

const (
	displayModeTariff1 uint8 = 1 << iota
	displayModeTariff2
	displayModeTariff3
	displayModeTariff4
	displayModeAmount
	displayModePower
	displayModeTime
	displayModeDate
)

const (
	displayModeTariffSchedule uint8 = 1 << iota
	displayModeUIF
	displayModeReactiveEnergy
	displayModeMaximumResets
	displayModeWorkingTime
	displayModeBatteryLifetime
	displayModePowerLimit
	displayModeEnergyLimit

	MaxEventsIndex = 0x3F
	CurrentMonth   = 0x0F

	MaximumPower    = 0x0
	MaximumAmperage = 0x1
	MaximumVoltage  = 0x2
)

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
	CommandWriteDisplayMode uint8 = 0x09 // установка режима индикации
	// 0x0A // установка числа действующих тарифов
	// 0x0B // установка тарифа
	// 0x0C // сброс защёлки "напряжение батареи"
	CommandWriteDisplayTime uint8 = 0x0D // установка времени индикации
	// 0x0F // установка режима лимита мощности
	// 0x10 // установка таблицы праздничных дней
	// 0x11 // установка таблицы переключений тарифных зон
	// 0x12 // сброс максимумов
	// 0x13 // сброс максимумов. Под перемычкой
	// 0x14 // сброс наработки батареи.

	// группа сетевых команд чтения
	CommandReadAddressGroup         = 0x20 // чтение группового адреса счетчика
	CommandReadDatetime             = 0x21 // чтение установленных даты и времени
	CommandReadPowerMaximum         = 0x22 // чтение лимита мощности
	CommandReadEnergyMaximum        = 0x23 // чтение лимита энергии
	CommandReadDaylightSavingTime   = 0x24 // чтение значения флага сезонного времени
	CommandReadTimeCorrection       = 0x25 // чтение величины коррекции времени
	CommandReadPowerCurrent         = 0x26 // чтение текущей мощности в нагрузке
	CommandReadPowerCounters        = 0x27 // чтение содержимого тарифных аккумуляторов
	CommandReadVersion              = 0x28 // чтение версии ПО
	CommandReadBatteryVoltage       = 0x29 // чтение напряжения встроенной батарейки
	CommandReadDisplayMode          = 0x2A // чтение режима индикации
	CommandReadLastPowerOffDatetime = 0x2B // чтение времени последнего отключения напряжения
	CommandReadLastPowerOnDatetime  = 0x2C // чтение времени последнего включения напряжения
	CommandReadOptocouplerFunction  = 0x2D // Чтение функции выходного оптрона
	CommandReadTariffCount          = 0x2E // чтение количества действующих тарифов
	CommandReadSerialNumber         = 0x2F // чтение серийного номера

	// группа дополнительных сетевых команд чтения
	CommandReadHolidays          = 0x30 // чтение таблицы праздничных дней
	CommandReadTariffZoneChanged = 0x31 // чтение таблицы переключений тарифных зон
	CommandReadMonthlyStat       = 0x32 // чтение месячных срезов
	CommandReadMaximum           = 0x33 // чтение месячных срезов
	CommandReadEventsPowerOnOff  = 0x34 // чтение буфера событий включений/выключения
	CommandReadEventsOpenClose   = 0x35 // чтение буфера событий включений/выключения
	CommandReadEventsParameters  = 0x36 // чтение буфера событий параметризации
	// 0x37 // чтение получасовых мощностей или суточных срезов
	// 0x38 // чтение месячных срезов реактивной энергии
	// 0x39 // чтение буфера событий качества электричества
	CommandReadEventsRelay = 0x3A // чтение буфера событий реле

	CommandReadCurrentTariff = 0x60 // чтение тарифа
	CommandReadLastOpenCap   = 0x61 // Чтение времени последнего вскрытия крышки счётчика
	CommandReadLastCloseCap  = 0x62 // чтение времение последнего закрытия крышки счетчика
	CommandReadParamsCurrent = 0x63 // чтение текущих параметров сети U, I, P
	// 0x64 // чтение коэффициента коррекции хода часов Введена для чтения коэффициента коррекции без перемычки.
	CommandReadModel           uint8 = 0x65 // чтение слова исполнения
	CommandReadMakeDate        uint8 = 0x66 // чтение даты изготовления
	CommandReadDisplayTime     uint8 = 0x67 // чтение времени индекации
	CommandEnergyLimitMode     uint8 = 0x68 // чтение режима лимита мощности
	CommandReadWorkingTime     uint8 = 0x69 // чтение времени наработки
	CommandReadDisplayModeExt  uint8 = 0x6A // чтение режима доп. индикации
	CommandReadParamLastChange uint8 = 0x6B // чтение времени последней парам. счётчика
	// 0x6C // чтение номера модема и системы и уровня сигнала
	CommandReadRelayMode                   = 0x6D // чтение режима управления реле
	CommandReadPowerLimits                 = 0x6E // чтение потарифных лимитов Энергии (остатки)
	CommandReadAllowIndicationUnderBattery = 0x6F // чтение флага разрешения индикации под батарейкой

	// Группа доп. сетевых команд записи
	//  = 0x70 // установка режима доп. индикации
	//  = 0x71 // установка режима управления реле
	//  = 0x72 // установка потарифных лимитов энергии
	//  = 0x73 // установка флага разрешения индикации под батарейкой
	//  = 0x74 // установка флага разрешения работы с модемом PLC2
	//  = 0x75 // установка множителя таймаута
	//  = 0x78 // установка тарифного расписания сжатым методом
	//  = 0x79 // обмен данными с PLC1 модемом
	//  = 0x7A // запись параметра

	// Группа доп. сетевых команд чтения
	//  = 0x80 // чтение флага разрешения модема PLC2
	//  = 0x81 // чтение доп. параметров сети (частота) и текущего тарифа
	//  = 0x82 // Чтение множителя таймаута
	//  = 0x85 // чтение содержимого тарифных аккумуляторов реактивной энергии
	//  = 0x86 // чтение параметра
)

type MercuryV1 struct {
	connection connection.Connection
	options    options
}

func New(conn connection.Connection, opts ...Option) *MercuryV1 {
	conn.ApplyOptions(connection.WithGlobalLock(true))
	conn.ApplyOptions(connection.WithOnceInit(true))

	m := &MercuryV1{
		connection: conn,
		options:    defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&m.options)
	}

	return m
}

func (m *MercuryV1) Invoke(request *Packet) (*Packet, error) {
	if !IsCommandSupported(m.options.device, request.Command()) {
		return nil, ErrCommandNotSupported
	}

	if request.Address() == 0 {
		request = request.WithAddress(m.options.address)
	}

	if request.Address() == 0 {
		return nil, errors.New("device address is empty")
	}

	requestData, err := request.MarshalBinary()
	if err != nil {
		return nil, err
	}

	responseData, err := m.connection.Invoke(requestData)
	if err != nil {
		return nil, err
	}

	response := NewPacket()

	if err = response.UnmarshalBinary(responseData); err != nil {
		return nil, err
	}

	// check ADDR
	if response.Address() != request.Address() {
		return nil, fmt.Errorf("error ADDR of response packet %X have %X want %X",
			responseData, response.Address(), request.Address())
	}

	return response, nil
}

package v1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

/*
	Чтение группового адреса счетчика

	CMD: 20h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-GADDR-CRC
*/
func (m *MercuryV1) AddressGroup() (address uint32, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadAddressGroup))
	if err == nil {
		address = response.PayloadAsBuffer().ReadUint32()
	}

	return address, nil
}

/*
	Чтение внутренних часов и календаря счетчика

	CMD: 21h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) Datetime() (date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadDatetime))
	if err == nil {
		date = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return date, err
}

/*
	Чтение лимита мощности

	CMD: 22h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-mpower-CRC
*/
func (m *MercuryV1) PowerMaximum() (maximum uint64, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadPowerMaximum))
	if err == nil {
		maximum = response.PayloadAsBuffer().ReadBCD(2) * 10
	}

	return maximum, err
}

/*
	Чтение лимита энергии за месяц

	CMD: 23h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-menerg-CRC
*/
func (m *MercuryV1) EnergyMaximum() (maximum uint64, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadEnergyMaximum))
	if err == nil {
		maximum = response.PayloadAsBuffer().ReadBCD(2) * 1000
	}

	return maximum, err
}

/*
	Чтение флага сезонного времени

	CMD: 24h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-flag-CRC
*/
func (m *MercuryV1) DaylightSavingTime() (flag bool, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadDaylightSavingTime))
	if err == nil {
		flag = response.PayloadAsBuffer().ReadBool()
	}

	return flag, err
}

/*
	Чтение величины коррекции времени

	CMD: 25h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timecor-CRC
*/
func (m *MercuryV1) TimeCorrection() (duration time.Duration, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadTimeCorrection))
	if err == nil {
		duration = time.Duration(response.PayloadAsBuffer().ReadUint8())
	}

	// Коррекция времени в счетчике Меркурий 200 возможна в диапазоне ±30 минут в течение года.
	// Коррекция времени в счетчике Меркурий 230 возможна в диапазоне ±30 секунд в течении суток.

	return duration, err
}

/*
	Чтение текущей мощности в нагрузке

	CMD: 26h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-m-CRC
*/
func (m *MercuryV1) PowerCurrent() (power uint64, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadPowerCurrent))
	if err == nil {
		power = response.PayloadAsBuffer().ReadBCD(2) * 10
	}

	return power, err
}

/*
	Чтение содержимого тарифных аккумуляторов активной энергии

	CMD: 27h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-count*4-CRC
*/
func (m *MercuryV1) PowerCounters() (values *TariffValues, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadPowerCounters))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		values = NewTariffValues(dataOut.ReadCount(), dataOut.ReadCount(), dataOut.ReadCount(), dataOut.ReadCount())
	}

	return values, err
}

/*
	Чтение содержимого тарифных аккумуляторов активной энергии

	CMD: 28h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-ver-DataVer-CRC
*/
func (m *MercuryV1) Version() (version string, date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadVersion))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		// FIXME: [1 0 0 6 4 21]
		// по доке 2 байта на версию 4 на дату, но такое впечатление что 3 и 3

		version = fmt.Sprintf("%d.%d.%d", dataOut.ReadUint8(), dataOut.ReadUint8(), dataOut.ReadUint8())
		date = dataOut.ReadDate()
	}

	return version, date, err
}

/*
	Чтение напряжения на литиевой батарее

	CMD: 29h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-VVVV-CRC
*/
func (m *MercuryV1) BatteryVoltage() (voltage float64, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadBatteryVoltage))
	if err == nil {
		voltage = float64(response.PayloadAsBuffer().ReadBCD(4)) / 100
	}

	return voltage, err
}

/*
	Чтение режима индикации

	CMD: 2Ah
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-displ-CRC
*/
func (m *MercuryV1) DisplayMode() (mode *DisplayMode, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadDisplayMode))
	if err == nil {
		mode = NewDisplayMode(response.PayloadAsBuffer().ReadUint8())
	}

	return mode, err
}

/*
	Чтение времени последнего отключения напряжения

	CMD: 2Bh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) LastPowerOffDatetime() (date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadLastPowerOffDatetime))
	if err == nil {
		date = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return date, err
}

/*
	Чтение времени последнего включения напряжения

	CMD: 2Ch
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) LastPowerOnDatetime() (date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadLastPowerOnDatetime))
	if err == nil {
		date = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return date, err
}

/*
	Чтение количества действующих тарифов

	CMD: 2Dh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-function-CRC

	Функциональное назначение выходного ключа импульсного выхода
		0 - телеметрический выход 5000 имп/кВт.ч
		1 - телеметрический выход 10000 имп/кВт.ч
		2 - выход частоты встроенного кварца поделенной на 8
		3 - управление нагрузкой
*/
func (m *MercuryV1) OptocouplerFunction() (count uint8, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadOptocouplerFunction))
	if err == nil {
		count = response.PayloadAsBuffer().ReadUint8()
	}

	return count, err
}

/*
	Чтение количества действующих тарифов

	CMD: 2Eh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-tarif-CRC
*/
func (m *MercuryV1) TariffCount() (count uint8, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadTariffCount))
	if err == nil {
		count = response.PayloadAsBuffer().ReadUint8()
	}

	return count, err
}

/*
	Чтение серийного номера

	CMD: 2Fh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-serial-CRC
*/
func (m *MercuryV1) SerialNumber() (sn uint32, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadSerialNumber))
	if err == nil {
		sn = response.PayloadAsBuffer().ReadUint32()
	}

	return sn, err
}

/*
	Чтение таблицы праздничных дней

	CMD: 30h
	Request: ADDR-CMD-ii1-CRC
	Response: ADDR-CMD-(dd-mon)*8-CRC
*/
func (m *MercuryV1) Holidays() ([]time.Time, error) {
	request := NewPacket().
		WithCommand(CommandReadHolidays).
		WithPayload([]byte{0})

	response1, err := m.Invoke(request)
	if err != nil {
		return nil, err
	}

	request = request.WithPayload([]byte{1})
	response2, err := m.Invoke(request)
	if err != nil {
		return nil, err
	}

	response := append(response1.Payload(), response2.Payload()...)
	days := make([]time.Time, 0)
	year := time.Now().Year()

	for i := 0; i < len(response); i += 2 {
		if response[i] < 1 || response[i] > 31 || response[i+1] > 12 || response[i+1] < 1 {
			continue
		}

		days = append(days, time.Date(year, time.Month(response[i+1]), int(response[i]), 0, 0, 0, 0, time.UTC))
	}

	return days, nil
}

/*
	Чтение таблицы переключений тарифных зон

	CMD: 31h
	Request: ADDR-CMD-ii2-CRC
	Response: ADDR-CMD-(nh-mm)*16-CRC

	nh - Часы временной точки смены тарифа. В двух старших битах заложен номер тарифа 00 – 1, 01 – 2, 10 – 3, 11 – 4.
*/
func (m *MercuryV1) tariffZoneChanged(month uint8) (zones [][]uint8, err error) {
	var response *Packet

	request := NewPacket().
		WithCommand(CommandReadTariffZoneChanged).
		WithPayload([]byte{month})

	response, err = m.Invoke(request)
	if err == nil {
		dataOut := response.PayloadAsBuffer()
		zones = make([][]uint8, 0)

		for dataOut.Len() > 0 {
			value := dataOut.ReadUint8()

			if value == 255 {
				break
			}

			tariffBit := value >> 6
			hourBit := value & ^(tariffBit << 6)

			hour, _ := strconv.ParseUint(hex.EncodeToString([]byte{hourBit}), 10, 64)
			minute := dataOut.ReadBCD(1)

			zones = append(zones, []uint8{tariffBit + 1, uint8(hour), uint8(minute)})
		}
	}

	return zones, err
}

func (m *MercuryV1) ReadTariffZoneChanged() ([][]uint8, error) {
	return m.tariffZoneChanged(CurrentMonth)
}

func (m *MercuryV1) ReadTariffZoneChangedByMonth(month time.Month) ([][]uint8, error) {
	switch month {
	case time.January, time.February, time.March, time.April, time.May, time.June,
		time.July, time.August, time.September, time.October, time.November, time.December:
	default:
		return nil, errors.New("wrong month " + strconv.FormatInt(int64(month), 16))
	}

	return m.tariffZoneChanged(uint8(month) - 1)
}

/*
	Чтение месячных срезов

	CMD: 32h
	Request: ADDR-CMD-ii3-CRC
	Response: ADDR-CMD-count*4-CRC
*/
func (m *MercuryV1) monthlyStat(month uint8) (values *TariffValues, err error) {
	request := NewPacket().
		WithCommand(CommandReadMonthlyStat).
		WithPayload([]byte{month})

	response, err := m.Invoke(request)
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		values = NewTariffValues(dataOut.ReadCount(), dataOut.ReadCount(), dataOut.ReadCount(), dataOut.ReadCount())
	}

	return values, err
}

// 0x0F текущий месяц, но модель 200 возвращает не корректные значения
// поэтому лучше указывать месяц явно
func (m *MercuryV1) MonthlyStat() (*TariffValues, error) {
	return m.monthlyStat(CurrentMonth)
}

// значения счетчика на 1 число месяца
func (m *MercuryV1) MonthlyStatByMonth(month time.Month) (*TariffValues, error) {
	switch month {
	case time.January, time.February, time.March, time.April, time.May, time.June,
		time.July, time.August, time.September, time.October, time.November, time.December:
	default:
		return nil, errors.New("wrong month " + strconv.FormatInt(int64(month), 16))
	}

	return m.monthlyStat(uint8(month) - 1)
}

/*
	Чтение максимумов

	CMD: 33h
	Request: ADDR-CMD-ii4-CRC
	Response: ADDR-CMD-max-maxr-CRC
*/
func (m *MercuryV1) maximum(option uint8) (max uint64, maxDate time.Time, maxReset uint64, maxResetDate time.Time, err error) {
	request := NewPacket().
		WithCommand(CommandReadMaximum).
		WithPayload([]byte{option})

	response, err := m.Invoke(request)
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		max = dataOut.ReadBCD(2)
		maxDate = dataOut.ReadTimeDate(m.options.location)
		maxReset = dataOut.ReadBCD(2)
		maxResetDate = dataOut.ReadTimeDate(m.options.location)
	}

	return max, maxDate, maxReset, maxResetDate, err
}

func (m *MercuryV1) MaximumPower() (power uint64, date time.Time, powerReset uint64, dateReset time.Time, err error) {
	return m.maximum(MaximumPower)
}

func (m *MercuryV1) MaximumAmperage() (amperage float64, maxDate time.Time, amperageReset float64, maxResetDate time.Time, err error) {
	max, maxDate, maxReset, maxResetDate, err := m.maximum(MaximumAmperage)
	if err == nil {
		amperage = float64(max) / 100
		amperageReset = float64(maxReset) / 100
	}

	return amperage, maxDate, amperageReset, maxResetDate, err
}

func (m *MercuryV1) MaximumVoltage() (voltage uint64, maxDate time.Time, voltageReset uint64, maxResetDate time.Time, err error) {
	voltage, maxDate, voltageReset, maxResetDate, err = m.maximum(MaximumVoltage)
	if err == nil {
		voltage /= 10
		voltageReset /= 10
	}

	return voltage, maxDate, voltageReset, maxResetDate, err
}

/*
	Чтение буфера событий вкл/выкл

	CMD: 34h
	Request: ADDR-CMD-ii5-CRC
	Response: ADDR-CMD-event1-CRC
*/
func (m *MercuryV1) EventsPowerOnOff(index uint8) (event bool, t time.Time, err error) {
	if index > MaxEventsIndex {
		err = errors.New("wrong index value #" + strconv.FormatUint(uint64(index), 16))
	} else {
		var response *Packet

		request := NewPacket().
			WithCommand(CommandReadEventsPowerOnOff).
			WithPayload([]byte{index})

		response, err = m.Invoke(request)
		if err == nil {
			dataOut := response.PayloadAsBuffer()

			event = !dataOut.ReadBool()
			t = dataOut.ReadTimeDate(m.options.location)
		}
	}

	return event, t, err
}

/*
	Чтение буфера событий отк/закр

	CMD: 35h
	Request: ADDR-CMD-ii5-CRC
	Response: ADDR-CMD-event2-CRC
*/
func (m *MercuryV1) EventsOpenClose(index uint8) (event bool, t time.Time, err error) {
	if index > MaxEventsIndex {
		err = errors.New("wrong index value #" + strconv.FormatUint(uint64(index), 16))
	} else {
		var response *Packet

		request := NewPacket().
			WithCommand(CommandReadEventsOpenClose).
			WithPayload([]byte{index})

		response, err = m.Invoke(request)
		if err == nil {
			dataOut := response.PayloadAsBuffer()

			event = !dataOut.ReadBool()
			t = dataOut.ReadTimeDate(m.options.location)
		}
	}

	return event, t, err
}

/*
	Чтение буфера событий параметризации

	CMD: 36h
	Request: ADDR-CMD-ii5-CRC
	Response: ADDR-CMD-event3-CRC

	event3: ev3-ev4-ev5-ev6-ev7-dd-mon-yy
*/
func (m *MercuryV1) EventsParameters(index uint8) (t time.Time, err error) {
	if index > MaxEventsIndex {
		err = errors.New("wrong index value #" + strconv.FormatUint(uint64(index), 16))
	} else {
		var response *Packet

		request := NewPacket().
			WithCommand(CommandReadEventsParameters).
			WithPayload([]byte{index})

		response, err = m.Invoke(request)
		if err == nil {
			// TODO:
			fmt.Println(response.Payload())
		}
	}

	return t, err
}

/*
	Чтение буфера событий реле

	CMD: 3Аh
	Request: ADDR-CMD-ii5-CRC
	Response: ADDR-CMD-event5-CRC
*/
func (m *MercuryV1) EventsRelay(index uint8) (tariff uint8, err error) {
	if index > MaxEventsIndex {
		err = errors.New("wrong index value #" + strconv.FormatUint(uint64(index), 16))
	} else {
		var response *Packet

		request := NewPacket().
			WithCommand(CommandReadEventsRelay).
			WithPayload([]byte{index})

		response, err = m.Invoke(request)
		if err == nil {
			dataOut := response.PayloadAsBuffer()

			// TODO:
			fmt.Println(response.payload)
			fmt.Println(dataOut.ReadTimeDate(m.options.location))
		}
	}

	return tariff, err
}

/*
	Чтение тарифа

	CMD: 60h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD- Tarif-CRC
*/
func (m *MercuryV1) CurrentTariff() (tariff uint8, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadCurrentTariff))
	if err == nil {
		tariff = response.PayloadAsBuffer().ReadUint8()
	}

	return tariff, err
}

/*
	Чтение времени последнего вскрытия крышки счётчика

	CMD: 61h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) LastOpenCap() (date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadLastOpenCap))
	if err == nil {
		date = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return date, err
}

/*
	Чтение времени последнего закрытия крышки счётчика

	CMD: 62h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) LastCloseCap() (date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadLastCloseCap))
	if err == nil {
		date = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return date, err
}

/*
	Чтение значений U,I,P

	CMD: 63h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-V-I-P-CRC
*/
func (m *MercuryV1) UIPCurrent() (voltage uint64, amperage float64, power uint64, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadParamsCurrent))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		voltage = dataOut.ReadBCD(2) / 10
		amperage = float64(dataOut.ReadBCD(2)) / 100
		power = dataOut.ReadBCD(3)
	}

	return voltage, amperage, power, err
}

/*
	Чтение слова исполнения

	CMD: 65h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-ispoln-CRC

	https://github.com/instalator/ioBroker.mercury/blob/fd1195fd4695c513ae4a12e211e0a3dd290f21dc/lib/mercury.js#L1294
*/
func (m *MercuryV1) Model() (twoSensors, relay bool, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadModel))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		twoSensors = dataOut.ReadUint8() > 0
		relay = dataOut.ReadUint8() > 0
	}

	return twoSensors, relay, err
}

/*
	Чтение даты изготовления

	CMD: 66h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-datefabric-CRC
*/
func (m *MercuryV1) MakeDate() (date time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadMakeDate))
	if err == nil {
		date = response.PayloadAsBuffer().ReadDate()
	}

	return date, err
}

/*
	Чтение времени индикации

	CMD: 67h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-TIMEDISPL-CRC
*/
func (m *MercuryV1) DisplayTime() (values *TariffValues, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadDisplayTime))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		values = NewTariffValues(
			uint64(dataOut.ReadUint8()),
			uint64(dataOut.ReadUint8()),
			uint64(dataOut.ReadUint8()),
			uint64(dataOut.ReadUint8()))
	}

	return values, err
}

/*
	Чтение времени наработки

	CMD: 68h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-TLM-TILM-CRC
*/
func (m *MercuryV1) EnergyLimitMode() (step uint8, without time.Duration, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandEnergyLimitMode))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		step = dataOut.ReadUint8()

		// TODO:
		fmt.Println(response.Payload())
	}

	return step, without, err
}

/*
	Чтение времени наработки

	CMD: 69h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-TL-TLB-CRC
*/
func (m *MercuryV1) WorkingTime() (under, without time.Duration, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadWorkingTime))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		under = time.Duration(dataOut.ReadBCD(3)) * time.Hour
		without = time.Duration(dataOut.ReadBCD(3)) * time.Hour
	}

	return under, without, err
}

/*
	Чтение режима доп. индикации

	CMD: 6Bh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-displ1-TIMED-CRC
*/
func (m *MercuryV1) DisplayModeExt() (mode *DisplayModeExt, timed uint8, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadDisplayModeExt))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		mode = NewDisplayModeExt(dataOut.ReadUint8())
		timed = dataOut.ReadUint8()
	}

	return mode, timed, err
}

/*
	Чтение времени последней парам. счётчика

	CMD: 6Bh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) ParamLastChange() (datetime time.Time, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadParamLastChange))
	if err == nil {
		datetime = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return datetime, err
}

/*
	Чтение режима управления реле

	CMD: 6Dh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-RELE-CRC
*/
func (m *MercuryV1) RelayMode() (byLimits, buttonEmulation, enabled bool, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadRelayMode))
	if err == nil {
		enabled = true

		switch response.Payload()[0] {
		case 0x55:
			byLimits = true
		case 0x5A:
			buttonEmulation = true
		case 0xAA:
			enabled = false
		}
	}

	return byLimits, buttonEmulation, enabled, err
}

/*
	Чтение потарифных лимитов Энергии (остатки)

	CMD: 6Eh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-mp1-mp2-mp3-mp4-CRC
*/
func (m *MercuryV1) PowerLimits() (values *TariffValues, flag1, flag2, flag3, flag4 bool, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadPowerLimits))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		t1 := dataOut.ReadBCD(3)
		flag1 = dataOut.ReadUint8() == 0xAA
		t2 := dataOut.ReadBCD(3)
		flag2 = dataOut.ReadUint8() == 0xAA
		t3 := dataOut.ReadBCD(3)
		flag3 = dataOut.ReadUint8() == 0xAA
		t4 := dataOut.ReadBCD(3)
		flag4 = dataOut.ReadUint8() == 0xAA

		values = NewTariffValues(t1, t2, t3, t4)
	}

	return values, flag1, flag2, flag3, flag4, err
}

/*
	Чтение флага разрешения индикации под батарейкой

	CMD: 6Fh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-flag1-CRC
*/
func (m *MercuryV1) AllowIndicationUnderBattery() (flag bool, err error) {
	response, err := m.Invoke(NewPacket().WithCommand(CommandReadAllowIndicationUnderBattery))
	if err == nil {
		flag = response.Payload()[0] == 0x55
	}

	return flag, err
}

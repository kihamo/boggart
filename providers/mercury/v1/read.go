package v1

import (
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
func (m *MercuryV1) AddressGroup() (address []byte, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadAddressGroup))
	if err == nil {
		address = response.PayloadAsBuffer().Next(4)
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
	response, err := m.Invoke(NewRequest(RequestCommandReadDatetime))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadPowerMaximum))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadEnergyMaximum))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadDaylightSavingTime))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadTimeCorrection))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadPowerCurrent))
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
func (m *MercuryV1) PowerCounters() (t1, t2, t3, t4 uint64, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadPowerCounters))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		t1 = dataOut.ReadCount()
		t2 = dataOut.ReadCount()
		t3 = dataOut.ReadCount()
		t4 = dataOut.ReadCount()
	}

	return t1, t2, t3, t4, err
}

/*
	Чтение содержимого тарифных аккумуляторов активной энергии

	CMD: 28h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-ver-DataVer-CRC
*/
func (m *MercuryV1) Version() (version string, date time.Time, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadVersion))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadBatteryVoltage))
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
func (m *MercuryV1) DisplayMode() (t1, t2, t3, t4, amount, power, t, date bool, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadDisplayMode))
	if err == nil {
		bit := response.PayloadAsBuffer().ReadUint8()
		t1, t2, t3, t4, amount, power, t, date = bit&displayModeTariff1 != 0, bit&displayModeTariff2 != 0,
			bit&displayModeTariff3 != 0, bit&displayModeTariff4 != 0,
			bit&displayModeAmount != 0, bit&displayModePower != 0,
			bit&displayModeTime != 0, bit&displayModeDate != 0
	}

	return t1, t2, t3, t4, amount, power, t, date, err
}

/*
	Чтение времени последнего отключения напряжения

	CMD: 2Bh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-timedate-CRC
*/
func (m *MercuryV1) LastPowerOffDatetime() (date time.Time, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadLastPowerOffDatetime))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadLastPowerOnDatetime))
	if err == nil {
		date = response.PayloadAsBuffer().ReadTimeDateWithDayOfWeek(m.options.location)
	}

	return date, err
}

/*
	Чтение серийного номера

	CMD: 2Fh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-serial-CRC
*/
func (m *MercuryV1) SerialNumber() (sn []byte, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadSerialNumber))
	if err == nil {
		sn = response.PayloadAsBuffer().Next(4)
	}

	// TODO: конвертировать в что-то человекочитаемое

	return sn, err
}

/*
	Чтение количества действующих тарифов

	CMD: 2Eh
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-tarif-CRC
*/
func (m *MercuryV1) TariffCount() (count uint8, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadTariffCount))
	if err == nil {
		count = response.PayloadAsBuffer().ReadUint8()
	}

	return count, err
}

/*
	Чтение таблицы праздничных дней

	CMD: 30h
	Request: ADDR-CMD-ii1-CRC
	Response: ADDR-CMD-(dd-mon)*8-CRC
*/
func (m *MercuryV1) Holidays() ([]time.Time, error) {
	request := NewRequest(RequestCommandReadHolidays).
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
	Чтение месячных срезов

	CMD: 32h
	Request: ADDR-CMD-ii3-CRC
	Response: ADDR-CMD-count*4 -CRC
*/
func (m *MercuryV1) monthlyStat(month byte) (t1, t2, t3, t4 uint64, err error) {
	request := NewRequest(RequestCommandReadMonthlyStat).
		WithPayload([]byte{month})

	response, err := m.Invoke(request)
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		t1 = dataOut.ReadCount()
		t2 = dataOut.ReadCount()
		t3 = dataOut.ReadCount()
		t4 = dataOut.ReadCount()
	}

	return t1, t2, t3, t4, err
}

// 0x0F текущий месяц, но модель 200 возвращает не корректные значения
// поэтому лучше указывать месяц явно
func (m *MercuryV1) MonthlyStat() (uint64, uint64, uint64, uint64, error) {
	return m.monthlyStat(0x0F)
}

// значения счетчика на 1 число месяца
func (m *MercuryV1) MonthlyStatByMonth(month time.Month) (uint64, uint64, uint64, uint64, error) {
	return m.monthlyStat(byte(int(month) - 1))
}

/*
	Чтение максимумов

	CMD: 33h
	Request: ADDR-CMD-ii4-CRC
	Response: ADDR-CMD-max-maxr-CRC
*/
func (m *MercuryV1) maximum(option uint8) (max uint64, maxDate time.Time, maxReset uint64, maxResetDate time.Time, err error) {
	request := NewRequest(RequestCommandReadMaximum).
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
		var response *Response

		request := NewRequest(RequestCommandReadEventsPowerOnOff).
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
		var response *Response

		request := NewRequest(RequestCommandReadEventsOpenClose).
			WithPayload([]byte{uint8(index)})

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
	Чтение тарифа

	CMD: 60h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD- Tarif-CRC
*/
func (m *MercuryV1) CurrentTariff() (tariff uint8, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadCurrentTariff))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadLastOpenCap))
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
	response, err := m.Invoke(NewRequest(RequestCommandReadLastCloseCap))
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
func (m *MercuryV1) ParamsCurrent() (voltage uint64, amperage float64, power uint64, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadParamsCurrent))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		voltage = dataOut.ReadBCD(2) / 10
		amperage = float64(dataOut.ReadBCD(2)) / 100
		power = dataOut.ReadBCD(3)
	}

	return voltage, amperage, power, err
}

/*
	Чтение даты изготовления

	CMD: 63h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-datefabric-CRC
*/
func (m *MercuryV1) MakeDate() (date time.Time, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadMakeDate))
	if err == nil {
		date = response.PayloadAsBuffer().ReadDate()
	}

	return date, err
}

/*
	Чтение даты изготовления

	CMD: 67h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-TIMEDISPL-CRC
*/
func (m *MercuryV1) DisplayTime() (t1, t2, t3, t4 uint8, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadDisplayTime))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		t1 = dataOut.ReadUint8()
		t2 = dataOut.ReadUint8()
		t3 = dataOut.ReadUint8()
		t4 = dataOut.ReadUint8()
	}

	return t1, t2, t3, t4, err
}

/*
	Чтение времени наработки

	CMD: 69h
	Request: ADDR-CMD-CRC
	Response: ADDR-CMD-TL-TLB-CRC
*/
func (m *MercuryV1) WorkingTime() (under, without time.Duration, err error) {
	response, err := m.Invoke(NewRequest(RequestCommandReadWorkingTime))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		under = time.Duration(dataOut.ReadBCD(3)) * time.Hour
		without = time.Duration(dataOut.ReadBCD(3)) * time.Hour
	}

	return
}

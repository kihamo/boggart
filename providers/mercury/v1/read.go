package v1

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (m *MercuryV1) AddressGroup() ([]byte, error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadAddressGroup,
	})

	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func (m *MercuryV1) Datetime() (date time.Time, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadDatetime,
	})

	if err == nil {
		date = ParseDatetime(response.Payload, m.options.location)
	}

	return
}

func (m *MercuryV1) SerialNumber() (sn int64, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadSerialNumber,
	})

	if err != nil {
		sn, err = strconv.ParseInt(hex.EncodeToString(response.Payload), 16, 0)
	}

	return
}

func (m *MercuryV1) WordType() error {
	response, err := m.Request(&Request{
		Command: RequestCommandReadWordType,
	})
	if err != nil {
		return err
	}

	// TODO:
	_ = response

	return nil
}

func (m *MercuryV1) MakeDate() (date time.Time, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadMakeDate,
	})

	if err == nil {
		date = ParseDate(response.Payload, m.options.location)
	}

	return
}

func (m *MercuryV1) Version() (version string, date time.Time, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadVersion,
	})

	if err == nil {
		version = fmt.Sprintf("%d.%d.%d", ParseInt(response.Payload[0]), ParseInt(response.Payload[1]), ParseInt(response.Payload[2]))
		date = ParseDate(response.Payload[3:], m.options.location)
	}

	return
}

// PowerMaximum return maximum of power in W
func (m *MercuryV1) PowerMaximum() (maximum int64, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadPowerMaximum,
	})

	if err == nil {
		maximum = ParseInt(response.Payload...) * 10
	}

	return
}

// EnergyMaximum return maximum of energy in W/h
func (m *MercuryV1) EnergyMaximum() (maximum int64, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadEnergyMaximum,
	})

	if err == nil {
		maximum = ParseInt(response.Payload...) * 1000
	}

	return
}

// BatteryVoltage return voltage of battery in V
func (m *MercuryV1) BatteryVoltage() (voltage float64, err error) {
	response, err := m.Request(&Request{
		Command: RequestCommandReadBatteryVoltage,
	})

	if err == nil {
		voltage = float64(ParseInt(response.Payload...)) / 100
	}

	return
}

func (m *MercuryV1) DisplayMode() (t1, t2, t3, t4, amount, power, t, date bool, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadDisplayMode,
	})

	if err == nil {
		bit := int(r.Payload[0])
		t1, t2, t3, t4, amount, power, t, date = bit&displayModeTariff1 != 0, bit&displayModeTariff2 != 0,
			bit&displayModeTariff3 != 0, bit&displayModeTariff4 != 0,
			bit&displayModeAmount != 0, bit&displayModePower != 0,
			bit&displayModeTime != 0, bit&displayModeDate != 0
	}

	return
}

// PowerCounters returns value of T1, T2, T3 and T4 in W/h
func (m *MercuryV1) PowerCounters() (t1, t2, t3, t4 uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadPowerCounters,
	})

	if err == nil {
		values := make([]uint64, 4)
		for i := 0; i < len(r.Payload); i += 4 {
			values[i/4] = ParseUint(r.Payload[i:i+4]...) * 10
		}

		t1, t2, t3, t4 = values[0], values[1], values[2], values[3]
	}

	return
}

// PowerUser return power in W
func (m *MercuryV1) PowerCurrent() (power uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadPowerCurrent,
	})

	if err == nil {
		power = ParseUint(r.Payload...)
	}

	return
}

func (m *MercuryV1) DaylightSavingTime() (flag bool, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadDaylightSavingTime,
	})

	if err == nil {
		flag = !bytes.Equal(r.Payload, []byte{0})
	}

	return
}

func (m *MercuryV1) TimeCorrection() (duration uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadTimeCorrection,
	})

	if err == nil {
		duration = uint64(r.Payload[0])
	}

	// Коррекция времени в счетчике Меркурий 200 возможна в диапазоне ±30 минут в течение года.
	// Коррекция времени в счетчике Меркурий 230 возможна в диапазоне ±30 секунд в течении суток.

	return
}

// ParamsCurrent returns current value of voltage in V, amperage in A, power in W
func (m *MercuryV1) ParamsCurrent() (voltage uint64, amperage float64, power uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadParamsCurrent,
	})

	if err == nil {
		voltage = ParseUint(r.Payload[0:2]...) / 10
		amperage = ParseFloat(r.Payload[2:4]...) / 100
		power = ParseUint(r.Payload[4:7]...)
	}

	return
}

func (m *MercuryV1) LastPowerOffDatetime() (date time.Time, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadLastPowerOffDatetime,
	})

	if err == nil {
		date = ParseDatetime(r.Payload, m.options.location)
	}

	return
}

func (m *MercuryV1) LastPowerOnDatetime() (date time.Time, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadLastPowerOnDatetime,
	})

	if err == nil {
		date = ParseDatetime(r.Payload, m.options.location)
	}

	return
}

func (m *MercuryV1) LastCloseCap() (date time.Time, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadLastCloseCap,
	})

	if err == nil {
		date = ParseDatetime(r.Payload, m.options.location)
	}

	return
}

func (m *MercuryV1) TariffCount() (count uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadTariffCount,
	})

	if err == nil {
		count = uint64(r.Payload[0])
	}

	return
}

func (m *MercuryV1) Holidays() ([]time.Time, error) {
	response1, err := m.Request(&Request{
		Command: RequestCommandReadHolidays,
		Payload: []byte{0},
	})
	if err != nil {
		return nil, err
	}

	response2, err := m.Request(&Request{
		Command: RequestCommandReadHolidays,
		Payload: []byte{1},
	})
	if err != nil {
		return nil, err
	}

	response := append(response1.Payload, response2.Payload...)
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

func (m *MercuryV1) MaximumPower() (power int64, date time.Time, powerReset int64, dateReset time.Time, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadMaximum,
		Payload: []byte{MaximumPower},
	})

	if err == nil {
		power = ParseInt(r.Payload[:2]...)
		date = ParseDatetime(r.Payload[2:8], m.options.location)
		powerReset = ParseInt(r.Payload[8:10]...)
		dateReset = ParseDatetime(r.Payload[10:], m.options.location)
	}

	return
}

func (m *MercuryV1) MaximumAmperage() (amperage float64, date time.Time, amperageReset float64, dateReset time.Time, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadMaximum,
		Payload: []byte{MaximumAmperage},
	})

	if err == nil {
		amperage = ParseFloat(r.Payload[:2]...) / 100
		date = ParseDatetime(r.Payload[2:8], m.options.location)
		amperageReset = ParseFloat(r.Payload[8:10]...) / 100
		dateReset = ParseDatetime(r.Payload[10:], m.options.location)
	}

	return
}

func (m *MercuryV1) MaximumVoltage() (voltage uint64, date time.Time, voltageReset uint64, dateReset time.Time, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadMaximum,
		Payload: []byte{MaximumVoltage},
	})

	if err == nil {
		voltage = ParseUint(r.Payload[:2]...) / 10
		date = ParseDatetime(r.Payload[2:8], m.options.location)
		voltageReset = ParseUint(r.Payload[8:10]...) / 10
		dateReset = ParseDatetime(r.Payload[10:], m.options.location)
	}

	return
}

func (m *MercuryV1) monthlyStat(month byte) (t1, t2, t3, t4 uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadMonthlyStat,
		Payload: []byte{month},
	})

	if err == nil {
		values := make([]uint64, 4)
		for i := 0; i < len(r.Payload); i += 4 {
			values[i/4] = ParseUint(r.Payload[i:i+4]...) * 10
		}

		t1, t2, t3, t4 = values[0], values[1], values[2], values[3]
	}

	return
}

func (m *MercuryV1) MonthlyStat() (uint64, uint64, uint64, uint64, error) {
	// 0x0F текущий месяц, но модель 200 возвращает не корректные значения
	// поэтому лучше указывать месяц явно

	return m.monthlyStat(0x0F)
}

// значения счетчика на 1 число месяца
func (m *MercuryV1) MonthlyStatByMonth(month time.Month) (uint64, uint64, uint64, uint64, error) {
	return m.monthlyStat(byte(int(month) - 1))
}

func (m *MercuryV1) EventsPowerOnOff(index uint64) (event bool, t time.Time, err error) {
	if index > MaxEventsIndex {
		err = errors.New("wrong index value #" + strconv.FormatUint(index, 16))
	} else {
		var r *Response

		r, err = m.Request(&Request{
			Command: RequestCommandReadEventsPowerOnOff,
			Payload: []byte{uint8(index)},
		})

		if err == nil && !bytes.Equal(r.Payload[4:], []byte{255, 255, 255}) {
			event = r.Payload[0] != 1
			t = ParseDatetime(r.Payload, m.options.location)
		}
	}

	return event, t, err
}

func (m *MercuryV1) EventsOpenClose(index uint64) (event bool, t time.Time, err error) {
	if index > MaxEventsIndex {
		err = errors.New("wrong index value #" + strconv.FormatUint(index, 16))
	} else {
		var r *Response

		r, err = m.Request(&Request{
			Command: RequestCommandReadEventsOpenClose,
			Payload: []byte{uint8(index)},
		})

		if err == nil && !bytes.Equal(r.Payload[4:], []byte{255, 255, 255}) {
			event = r.Payload[0] == 0
			t = ParseDatetime(r.Payload, m.options.location)
		}
	}

	return event, t, err
}

func (m *MercuryV1) CurrentTariff() (tariff uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadCurrentTariff,
	})

	if err == nil {
		tariff = uint64(r.Payload[0])
	}

	return
}

func (m *MercuryV1) DisplayTime() (t1, t2, t3, t4 uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadDisplayTime,
	})

	if err == nil {
		t1, t2, t3, t4 = uint64(r.Payload[0]), uint64(r.Payload[1]), uint64(r.Payload[2]), uint64(r.Payload[3])
	}

	return
}

func (m *MercuryV1) WorkingTime() (under, without uint64, err error) {
	r, err := m.Request(&Request{
		Command: RequestCommandReadWorkingTime,
	})

	if err == nil {
		under = ParseUint(r.Payload[0:3]...)
		without = ParseUint(r.Payload[3:6]...)
	}

	return
}

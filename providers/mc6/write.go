package mc6

import (
	"errors"
)

// FIXME: на устройстве не срабатывает
func (m *MC6) TemperatureFormat(format uint16) error {
	return m.Write(AddressTemperatureFormat, format)
}

func (m *MC6) Status(flag bool) error {
	var value uint16

	if flag {
		value = 1
	}

	return m.Write(AddressStatus, value)
}

func (m *MC6) SetTemperature(value float64) error {
	value *= 10

	if value < 50 || value > 350 {
		return errors.New("wrong temperature value 50 >= value <= 350")
	}

	return m.Write(AddressSetTemperature, uint16(value))
}

func (m *MC6) Away(flag bool) error {
	var value uint16

	if flag {
		value = 1
	}

	return m.Write(AddressAway, value)
}

func (m *MC6) AwayTemperature(value uint16) error {
	value *= 10

	if value < 50 || value > 350 {
		return errors.New("wrong away temperature value 50 >= value <= 350")
	}

	return m.Write(AddressAwayTemperature, value)
}

func (m *MC6) HoldingTemperature(value uint16) error {
	value *= 10

	if value < 50 || value > 350 {
		return errors.New("wrong holding temperature value 50 >= value <= 350")
	}

	return m.Write(AddressHoldingTemperature, value)
}

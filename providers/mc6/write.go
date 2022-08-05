package mc6

import (
	"errors"
)

func (m *MC6) SetTemperature(value float64) error {
	value *= 10

	if value < 50 || value > 350 {
		return errors.New("wrong temperature value 50 >= value <= 350")
	}

	return m.Write(AddressSetTemperature, uint16(value))
}

package mc6

import (
	"fmt"
)

func (m *MC6) RoomTemperature() (float64, error) {
	value, err := m.Read(AddressRoomTemperature)
	if err != nil {
		return 0, err
	}

	// если датчик подключен не правильно, возвращается 999
	if value > 500 {
		return 0, fmt.Errorf("room sensor returned wrong value %d", value)
	}

	return float64(value) / 10, err
}

func (m *MC6) FloorTemperature() (float64, error) {
	value, err := m.Read(AddressFloorTemperature)
	if err != nil {
		return 0, err
	}

	// если датчик подключен не правильно, возвращается 999
	if value > 500 {
		return 0, fmt.Errorf("floor sensor returned wrong value %d", value)
	}

	return float64(value) / 10, err
}

func (m *MC6) Humidity() (uint16, error) {
	value, err := m.Read(AddressHumidity)
	if err != nil {
		return 0, err
	}

	if value > 99 {
		return 0, fmt.Errorf("floor sensor returned wrong value %d", value)
	}

	return value, err
}

// статус активации режима нагрева (реле замкнуто)
func (m *MC6) HeatingOutputStatus() (bool, error) {
	value, err := m.Read(AddressHeatingOutputStatus)
	if err != nil {
		return false, err
	}

	return value == 1, err
}

func (m *MC6) HoldingFunction() (bool, error) {
	value, err := m.Read(AddressHoldingFunction)
	if err != nil {
		return false, err
	}

	return value == 1, err
}

func (m *MC6) FloorOverheat() (bool, error) {
	value, err := m.Read(AddressFloorOverheat)
	if err != nil {
		return false, err
	}

	return value == 1, err
}

func (m *MC6) DeviceType() (uint16, error) {
	return m.Read(AddressDeviceType)
}

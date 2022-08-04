package mc6

import (
	"encoding/binary"
)

func (m *MC6) RoomTemperature() (float64, error) {
	response, err := m.Invoke(AddressRoomTemperature)
	if err != nil {
		return 0, err
	}

	value := binary.BigEndian.Uint16(response)

	return float64(value) / 10, err
}

func (m *MC6) FloorTemperature() (float64, error) {
	response, err := m.Invoke(AddressFloorTemperature)
	if err != nil {
		return 0, err
	}

	value := binary.BigEndian.Uint16(response)

	return float64(value) / 10, err
}

func (m *MC6) Humidity() (uint16, error) {
	response, err := m.Invoke(AddressHumidity)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(response), err
}

// статус активации режима нагрева (реле замкнуто)
func (m *MC6) HeatingOutputStatus() (bool, error) {
	response, err := m.Invoke(AddressHeatingOutputStatus)
	if err != nil {
		return false, err
	}

	return binary.BigEndian.Uint16(response) == 1, err
}

func (m *MC6) DeviceType() (uint16, error) {
	response, err := m.Invoke(AddressDeviceType)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(response), err
}

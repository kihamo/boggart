package mc6

func (m *MC6) RoomTemperature() (float64, error) {
	value, err := m.Read(AddressRoomTemperature)
	if err != nil {
		return 0, err
	}

	return float64(value) / 10, err
}

func (m *MC6) FloorTemperature() (float64, error) {
	value, err := m.Read(AddressFloorTemperature)
	if err != nil {
		return 0, err
	}

	return float64(value) / 10, err
}

func (m *MC6) Humidity() (uint16, error) {
	return m.Read(AddressHumidity)
}

// статус активации режима нагрева (реле замкнуто)
func (m *MC6) HeatingOutputStatus() (bool, error) {
	value, err := m.Read(AddressHeatingOutputStatus)
	if err != nil {
		return false, err
	}

	return value == 1, err
}

func (m *MC6) DeviceType() (uint16, error) {
	return m.Read(AddressDeviceType)
}

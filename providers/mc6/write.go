package mc6

// FIXME:
// HA на устройстве не срабатывает
// FCU4 любое значение вызывает перезагрузку устройства
func (m *MC6) SetTemperatureFormat(format uint16) error {
	return m.Write(AddressTemperatureFormat, 2, format)
}

func (m *MC6) SetStatus(flag bool) error {
	return m.WriteBool(AddressStatus, flag)
}

func (m *MC6) SetSystemMode(value uint16) error {
	return m.Write(AddressSystemMode, 1, value)
}

func (m *MC6) SetFanSpeed(value uint16) error {
	return m.Write(AddressFanSpeed, 1, value)
}

func (m *MC6) SetTargetTemperature(value float64) error {
	return m.WriteTemperature(AddressTargetTemperature, value)
}

func (m *MC6) SetAway(flag bool) error {
	return m.WriteBool(AddressAway, flag)
}

func (m *MC6) SetAwayTemperature(value float64) error {
	return m.WriteTemperature(AddressAwayTemperature, value)
}

func (m *MC6) SetHoldingTemperature(value float64) error {
	return m.WriteTemperature(AddressHoldingTemperature, value)
}

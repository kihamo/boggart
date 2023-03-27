package mc6

import (
	"time"
)

// FIXME:
// HA на устройстве не срабатывает
// FCU4 любое значение вызывает перезагрузку устройства
func (m *MC6) SetTemperatureFormat(format uint16) error {
	return m.WriteUint16(AddressTemperatureFormat, format)
}

func (m *MC6) SetStatus(flag bool) error {
	return m.WriteBool(AddressStatus, flag)
}

func (m *MC6) SetSystemMode(value uint16) error {
	return m.WriteUint16(AddressSystemMode, value)
}

func (m *MC6) SetFanSpeed(value uint16) error {
	return m.WriteUint16(AddressFanSpeed, value)
}

func (m *MC6) SetTargetTemperature(value float64) error {
	return m.WriteTemperature(AddressTargetTemperature, value)
}

func (m *MC6) SetAway(flag bool) error {
	return m.WriteBool(AddressAway, flag)
}

func (m *MC6) SetAwayTemperature(value uint16) error {
	return m.WriteTemperature(AddressAwayTemperature, float64(value))
}

func (m *MC6) SetHoldingTime(value time.Duration) error {
	return m.WriteUint32(AddressHoldingTime, uint32(value.Minutes()))
}

func (m *MC6) SetHoldingTemperature(value uint16) error {
	return m.WriteTemperature(AddressHoldingTemperature, float64(value))
}

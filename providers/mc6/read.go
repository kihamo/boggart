package mc6

import (
	"fmt"
)

func (m *MC6) RoomTemperature() (float64, error) {
	return m.ReadTemperature(AddressRoomTemperature)
}

func (m *MC6) FloorTemperature() (float64, error) {
	return m.ReadTemperature(AddressFloorTemperature)
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

// реле между HV - L замкнуто
func (m *MC6) HeatingValve() (bool, error) {
	return m.ReadBool(AddressHeatingValve)
}

// реле между СV - L замкнуто
func (m *MC6) CoolingValve() (bool, error) {
	return m.ReadBool(AddressCoolingValve)
}

// статус активации режима нагрева (реле замкнуто)
func (m *MC6) HeatingOutput() (bool, error) {
	return m.ReadBool(AddressHeatingOutput)
}

func (m *MC6) HoldingFunction() (bool, error) {
	return m.ReadBool(AddressHoldingFunction)
}

func (m *MC6) FloorOverheat() (bool, error) {
	return m.ReadBool(AddressFloorOverheat)
}

func (m *MC6) DeviceType() (Device, error) {
	value, err := m.Read(AddressDeviceType)

	if err != nil {
		return 0, err
	}

	return Device(value), err
}

func (m *MC6) TemperatureFormat() (uint16, error) {
	return m.Read(AddressTemperatureFormat)
}

func (m *MC6) Status() (bool, error) {
	return m.ReadBool(AddressStatus)
}

func (m *MC6) SystemMode() (uint16, error) {
	return m.Read(AddressSystemMode)
}

func (m *MC6) FanSpeed() (uint16, error) {
	return m.Read(AddressFanSpeed)
}

func (m *MC6) TargetTemperature() (float64, error) {
	return m.ReadTemperature(AddressTargetTemperature)
}

func (m *MC6) Away() (bool, error) {
	return m.ReadBool(AddressAwayTemperature)
}

func (m *MC6) AwayTemperature() (float64, error) {
	return m.ReadTemperature(AddressAwayTemperature)
}

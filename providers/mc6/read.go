package mc6

import (
	"fmt"
	"time"
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

// по факту не работает, всегда отдает 0 (возможно из-за невключения самой функции)
func (m *MC6) WindowsOpen() (bool, error) {
	return m.ReadBool(AddressWindowsOpen)
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

func (m *MC6) FanSpeedNumbers() (uint16, error) {
	value, err := m.Read(AddressFanSpeedNumbers)

	if err == nil {
		switch value {
		case 0:
			return 1, nil
		case 1:
			return 3, nil
		}
	}

	return value, err
}

// FIXME: по факту не работает, на HA всегда 80 на FCU всегда 0
func (m *MC6) TemperatureFormat() (uint16, error) {
	value, err := m.Read(AddressTemperatureFormat)

	// HA always return 80 for C
	if err == nil && value != 1 {
		return 0, err
	}

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
	return m.ReadBool(AddressAway)
}

func (m *MC6) AwayTemperature() (uint16, error) {
	return m.ReadTemperatureUint(AddressAwayTemperature)
}

func (m *MC6) HoldingTimeHi() (time.Duration, error) {
	return m.ReadDuration(AddressHoldingTimeHi)
}

func (m *MC6) HoldingTimeLow() (time.Duration, error) {
	return m.ReadDuration(AddressHoldingTimeLow)
}

func (m *MC6) HoldingTemperature() (uint16, error) {
	return m.ReadTemperatureUint(AddressHoldingTemperature)
}

func (m *MC6) PanelLock() (bool, error) {
	return m.ReadBool(AddressPanelLock)
}

func (m *MC6) PanelLockPin1() (uint16, error) {
	return m.Read(AddressPanelLockPin1)
}

func (m *MC6) PanelLockPin2() (uint16, error) {
	return m.Read(AddressPanelLockPin2)
}

func (m *MC6) PanelLockPin3() (uint16, error) {
	return m.Read(AddressPanelLockPin3)
}

func (m *MC6) PanelLockPin4() (uint16, error) {
	return m.Read(AddressPanelLockPin4)
}

func (m *MC6) TargetTemperatureMaximum() (uint16, error) {
	return m.ReadTemperatureUint(AddressTargetTemperatureMaximum)
}

func (m *MC6) TargetTemperatureMinimum() (uint16, error) {
	return m.ReadTemperatureUint(AddressTargetTemperatureMinimum)
}

func (m *MC6) FloorTemperatureLimit() (uint16, error) {
	return m.ReadTemperatureUint(AddressFloorTemperatureLimit)
}

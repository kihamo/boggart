package mc6

import (
	"encoding/binary"
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
	value, err := m.ReadUint16(AddressHumidity)
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
	value, err := m.ReadUint16(AddressDeviceType)

	if err != nil {
		return 0, err
	}

	return Device(value), err
}

func (m *MC6) FanSpeedNumbers() (uint16, error) {
	value, err := m.ReadUint16(AddressFanSpeedNumbers)

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
	value, err := m.ReadUint16(AddressTemperatureFormat)

	// HA always return 80 for C
	if err == nil && value != 1 {
		return 0, err
	}

	return value, err
}

func (m *MC6) Status() (bool, error) {
	return m.ReadBool(AddressStatus)
}

func (m *MC6) SystemMode() (uint16, error) {
	return m.ReadUint16(AddressSystemMode)
}

func (m *MC6) FanSpeed() (uint16, error) {
	return m.ReadUint16(AddressFanSpeed)
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

func (m *MC6) HoldingTime() (time.Duration, error) {
	return m.ReadDuration(AddressHoldingTime)
}

func (m *MC6) HoldingTemperatureAndTime() (float64, time.Duration, error) {
	response, err := m.Read(AddressHoldingTemperatureAndTime, 2)
	if err != nil {
		return 0, 0, err
	}

	tim := time.Duration(binary.BigEndian.Uint16(response[:2])) * time.Minute
	temperature := float64(binary.BigEndian.Uint16(response[2:])) / 10

	// TODO: validate temperature value

	return temperature, tim, err
}

func (m *MC6) HoldingTemperature() (uint16, error) {
	return m.ReadTemperatureUint(AddressHoldingTemperature)
}

func (m *MC6) PanelLock() (bool, error) {
	return m.ReadBool(AddressPanelLock)
}

func (m *MC6) PanelLockPin1() (uint16, error) {
	return m.ReadUint16(AddressPanelLockPin1)
}

func (m *MC6) PanelLockPin2() (uint16, error) {
	return m.ReadUint16(AddressPanelLockPin2)
}

func (m *MC6) PanelLockPin3() (uint16, error) {
	return m.ReadUint16(AddressPanelLockPin3)
}

func (m *MC6) PanelLockPin4() (uint16, error) {
	return m.ReadUint16(AddressPanelLockPin4)
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

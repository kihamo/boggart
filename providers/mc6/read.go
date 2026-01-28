package mc6

import (
	"encoding/binary"
	"fmt"
	"time"
)

// показания второго датчика (правый верхний угол на экране меньшими цифрами)
// может быть только датчиком пола (floor probe)
func (m *MC6) RoomTemperature() (float64, error) {
	return m.ReadTemperature(AddressRoomTemperature)
}

// показания первого датчика (в центре экрана большими цифрами)
// может быть одним из:
// - встроенный датчик
// - доп. датчик воздуха (remote air probe)
// - только датчик пола (floor probe)
func (m *MC6) FloorTemperature() (float64, error) {
	return m.ReadTemperature(AddressFloorTemperature)
}

func (m *MC6) Humidity() (uint16, error) {
	value, err := m.client.ReadHoldingRegistersUint16(AddressHumidity)
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
	return m.client.ReadHoldingRegistersBool(AddressHeatingValve)
}

// реле между СV - L замкнуто
func (m *MC6) CoolingValve() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressCoolingValve)
}

func (m *MC6) FanHigh() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressFanHigh)
}

func (m *MC6) FanMedium() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressFanMedium)
}

func (m *MC6) FanLow() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressFanLow)
}

// статус активации режима нагрева (реле замкнуто)
func (m *MC6) HeatingOutput() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressHeatingOutput)
}

func (m *MC6) Heat() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressHeat)
}

func (m *MC6) HotWater() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressHotWater)
}

func (m *MC6) TouchLock() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressTouchLock)
}

// по факту не работает, всегда отдает 0 (возможно из-за невключения самой функции)
func (m *MC6) WindowsOpen() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressWindowsOpen)
}

func (m *MC6) HolidayFunction() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressHolidayFunction)
}

func (m *MC6) HoldingFunction() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressHoldingFunction)
}

func (m *MC6) FloorOverheat() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressFloorOverheat)
}

func (m *MC6) DeviceType() (Device, error) {
	value, err := m.client.ReadHoldingRegistersUint16(AddressDeviceType)

	if err != nil {
		return 0, err
	}

	return Device(value), err
}

func (m *MC6) FanSpeedNumbers() (uint16, error) {
	value, err := m.client.ReadHoldingRegistersUint16(AddressFanSpeedNumbers)

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

func (m *MC6) SystemError() (bool, error) {
	value, err := m.client.ReadHoldingRegistersUint16(AddressSystemError)
	if err != nil {
		return false, err
	}

	// На HA возвращает 2 а не 1 в случае ошибки
	return value > 0, err
}

// FIXME: по факту не работает, на HA всегда 80 на FCU всегда 0
func (m *MC6) TemperatureFormat() (uint16, error) {
	value, err := m.client.ReadHoldingRegistersUint16(AddressTemperatureFormat)

	// HA always return 80 for C
	if err == nil && value != 1 {
		return 0, err
	}

	return value, err
}

func (m *MC6) Status() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressStatus)
}

func (m *MC6) SystemMode() (uint16, error) {
	return m.client.ReadHoldingRegistersUint16(AddressSystemMode)
}

func (m *MC6) FanSpeed() (uint16, error) {
	return m.client.ReadHoldingRegistersUint16(AddressFanSpeed)
}

func (m *MC6) TargetTemperature() (float64, error) {
	return m.ReadTemperature(AddressTargetTemperature)
}

func (m *MC6) Away() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressAway)
}

func (m *MC6) AwayTemperature() (uint16, error) {
	return m.ReadTemperatureUint(AddressAwayTemperature)
}

func (m *MC6) HoldingTime() (time.Duration, error) {
	return m.ReadDuration(AddressHoldingTime)
}

func (m *MC6) HoldingTemperatureAndTime() (float64, time.Duration, error) {
	response, err := m.client.ReadHoldingRegisters(AddressHoldingTemperatureAndTime, 2)
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

func (m *MC6) OptimumStart() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressOptimumStart)
}

func (m *MC6) Boost() (bool, error) {
        return m.client.ReadHoldingRegistersBool(AddressBoost)
}

func (m *MC6) PanelLock() (bool, error) {
	return m.client.ReadHoldingRegistersBool(AddressPanelLock)
}

func (m *MC6) PanelLockPin1() (uint16, error) {
	return m.client.ReadHoldingRegistersUint16(AddressPanelLockPin1)
}

func (m *MC6) PanelLockPin2() (uint16, error) {
	return m.client.ReadHoldingRegistersUint16(AddressPanelLockPin2)
}

func (m *MC6) PanelLockPin3() (uint16, error) {
	return m.client.ReadHoldingRegistersUint16(AddressPanelLockPin3)
}

func (m *MC6) PanelLockPin4() (uint16, error) {
	return m.client.ReadHoldingRegistersUint16(AddressPanelLockPin4)
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

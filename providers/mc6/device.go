package mc6

const (
	DeviceTypeHotWater             uint16 = 2
	DeviceTypeElectricHeating      uint16 = 3
	DeviceTypeFCU2                 uint16 = 4
	DeviceTypeFCU4                 uint16 = 5  // fcu-4, управление отоплением
	DeviceTypeHA                   uint16 = 30 // базовый простой MC6-HA, без горячей воды
	DeviceTypeElectricHeatingTimer uint16 = 31
)

type Device uint16

func NewDevice(t uint16) Device {
	return Device(t)
}

func (d Device) IsHotWater() bool {
	return uint16(d) == DeviceTypeHotWater
}

func (d Device) IsElectricHeating() bool {
	return uint16(d) == DeviceTypeElectricHeating
}

func (d Device) IsFCU2() bool {
	return uint16(d) == DeviceTypeFCU2
}

func (d Device) IsFCU4() bool {
	return uint16(d) == DeviceTypeFCU4
}

func (d Device) IsHA() bool {
	return uint16(d) == DeviceTypeHA
}

func (d Device) IsElectricHeatingTimer() bool {
	return uint16(d) == DeviceTypeElectricHeatingTimer
}

func (d Device) IsSupported(address uint16) bool {
	switch address {
	case AddressRoomTemperature, AddressHumidity, AddressTouchLock, AddressWindowsOpen, AddressHolidayFunction,
		AddressHoldingFunction, AddressBoostFunction, AddressDeviceType, AddressSystemError, AddressTemperatureFormat,
		AddressStatus, AddressSystemMode, AddressTargetTemperature, AddressAway, AddressAwayTemperature, AddressHoldingTime,
		AddressHoldingTemperatureAndTime, AddressHoldingTemperature, AddressHolidayStartTimeHigh, AddressHolidayStartTimeLow,
		AddressHolidayEndTimeHigh, AddressHolidayEndTimeLow, AddressBoostEndTimeHigh, AddressBoostEndTimeLow, AddressBoost,
		AddressPanelLock, AddressPanelLockPin1, AddressPanelLockPin2, AddressPanelLockPin3, AddressPanelLockPin4,
		AddressTargetTemperatureMaximum, AddressTargetTemperatureMinimum, AddressScheduleMode:
		return true
	case AddressFloorTemperature, AddressHeatingOutput, AddressFloorOverheat, AddressFloorTemperatureLimit:
		return d.IsHA()
	case AddressHeatingValve, AddressCoolingValve, AddressFanHigh, AddressFanMedium, AddressFanLow, AddressFanSpeedNumbers,
		AddressFanSpeed:
		return d.IsFCU4()

	case AddressValve, AddressHeat, AddressHotWater, AddressAuxiliaryHeat, AddressOptimumStart:
		return true // ????
	}

	return false
}

func (d Device) IsSupportedRoomTemperature() bool {
	return d.IsSupported(AddressRoomTemperature)
}

func (d Device) IsSupportedFloorTemperature() bool {
	return d.IsSupported(AddressFloorTemperature)
}

func (d Device) IsSupportedHumidity() bool {
	return d.IsSupported(AddressHumidity)
}

func (d Device) IsSupportedHeatingValve() bool {
	return d.IsSupported(AddressHeatingValve)
}

func (d Device) IsSupportedCoolingValve() bool {
	return d.IsSupported(AddressCoolingValve)
}

func (d Device) IsSupportedValve() bool {
	return d.IsSupported(AddressValve)
}

func (d Device) IsSupportedFanHigh() bool {
	return d.IsSupported(AddressFanHigh)
}

func (d Device) IsSupportedFanMedium() bool {
	return d.IsSupported(AddressFanMedium)
}

func (d Device) IsSupportedFanLow() bool {
	return d.IsSupported(AddressFanLow)
}

func (d Device) IsSupportedHeatingOutput() bool {
	return d.IsSupported(AddressHeatingOutput)
}

func (d Device) IsSupportedHeat() bool {
	return d.IsSupported(AddressHeat)
}

func (d Device) IsSupportedHotWater() bool {
	return d.IsSupported(AddressHotWater)
}

func (d Device) IsSupportedTouchLock() bool {
	return d.IsSupported(AddressTouchLock)
}

func (d Device) IsSupportedWindowsOpen() bool {
	return d.IsSupported(AddressWindowsOpen)
}

func (d Device) IsSupportedHolidayFunction() bool {
	return d.IsSupported(AddressHolidayFunction)
}

func (d Device) IsSupportedHoldingFunction() bool {
	return d.IsSupported(AddressHoldingFunction)
}

func (d Device) IsSupportedBoostFunction() bool {
	return d.IsSupported(AddressBoostFunction)
}

func (d Device) IsSupportedFloorOverheat() bool {
	return d.IsSupported(AddressFloorOverheat)
}

func (d Device) IsSupportedDeviceType() bool {
	return d.IsSupported(AddressDeviceType)
}

func (d Device) IsSupportedAuxiliaryHeat() bool {
	return d.IsSupported(AddressAuxiliaryHeat)
}

func (d Device) IsSupportedFanSpeedNumbers() bool {
	return d.IsSupported(AddressFanSpeedNumbers)
}

func (d Device) IsSupportedSystemError() bool {
	return d.IsSupported(AddressSystemError)
}

func (d Device) IsSupportedTemperatureFormat() bool {
	return d.IsSupported(AddressTemperatureFormat)
}

func (d Device) IsSupportedStatus() bool {
	return d.IsSupported(AddressStatus)
}

func (d Device) IsSupportedSystemMode() bool {
	return d.IsSupported(AddressSystemMode)
}

func (d Device) IsSupportedFanSpeed() bool {
	return d.IsSupported(AddressFanSpeed)
}

func (d Device) IsSupportedTargetTemperature() bool {
	return d.IsSupported(AddressTargetTemperature)
}

func (d Device) IsSupportedAway() bool {
	return d.IsSupported(AddressAway)
}

func (d Device) IsSupportedAwayTemperature() bool {
	return d.IsSupported(AddressAwayTemperature)
}

func (d Device) IsSupportedHoldingTime() bool {
	return d.IsSupported(AddressHoldingTime)
}

func (d Device) IsSupportedHoldingTemperatureAndTime() bool {
	return d.IsSupported(AddressHoldingTemperatureAndTime)
}

func (d Device) IsSupportedHoldingTemperature() bool {
	return d.IsSupported(AddressHoldingTemperature)
}

func (d Device) IsSupportedHolidayStartTimeHigh() bool {
	return d.IsSupported(AddressHolidayStartTimeHigh)
}

func (d Device) IsSupportedHolidayStartTimeLow() bool {
	return d.IsSupported(AddressHolidayStartTimeLow)
}

func (d Device) IsSupportedHolidayEndTimeHigh() bool {
	return d.IsSupported(AddressHolidayEndTimeHigh)
}

func (d Device) IsSupportedHolidayEndTimeLow() bool {
	return d.IsSupported(AddressHolidayEndTimeLow)
}

func (d Device) IsSupportedOptimumStart() bool {
	return d.IsSupported(AddressOptimumStart)
}

func (d Device) IsSupportedBoostEndTimeHigh() bool {
	return d.IsSupported(AddressBoostEndTimeHigh)
}

func (d Device) IsSupportedBoostEndTimeLow() bool {
	return d.IsSupported(AddressBoostEndTimeLow)
}

func (d Device) IsSupportedBoost() bool {
	return d.IsSupported(AddressBoost)
}

func (d Device) IsSupportedPanelLock() bool {
	return d.IsSupported(AddressPanelLock)
}

func (d Device) IsSupportedPanelLockPin1() bool {
	return d.IsSupported(AddressPanelLockPin1)
}

func (d Device) IsSupportedPanelLockPin2() bool {
	return d.IsSupported(AddressPanelLockPin2)
}

func (d Device) IsSupportedPanelLockPin3() bool {
	return d.IsSupported(AddressPanelLockPin3)
}

func (d Device) IsSupportedPanelLockPin4() bool {
	return d.IsSupported(AddressPanelLockPin4)
}

func (d Device) IsSupportedTargetTemperatureMaximum() bool {
	return d.IsSupported(AddressTargetTemperatureMaximum)
}

func (d Device) IsSupportedTargetTemperatureMinimum() bool {
	return d.IsSupported(AddressTargetTemperatureMinimum)
}

func (d Device) IsSupportedFloorTemperatureLimit() bool {
	return d.IsSupported(AddressFloorTemperatureLimit)
}

func (d Device) IsSupportedScheduleMode() bool {
	return d.IsSupported(AddressScheduleMode)
}

func (d Device) String() string {
	switch uint16(d) {
	case DeviceTypeHotWater:
		return "hot water"
	case DeviceTypeElectricHeating:
		return "electric heating"
	case DeviceTypeFCU2:
		return "FCU2"
	case DeviceTypeFCU4:
		return "FCU4"
	case DeviceTypeHA:
		return "HA"
	case DeviceTypeElectricHeatingTimer:
		return "electric heating timer"
	}

	return "unknown"
}

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
	case AddressRoomTemperature, AddressHumidity, AddressDeviceType, AddressTemperatureFormat, AddressStatus,
		AddressTargetTemperature, AddressTargetTemperatureMaximum, AddressTargetTemperatureMinimum:
		return true
	case AddressFloorTemperature, AddressFloorOverheat, AddressFloorTemperatureLimit:
		return d.IsHA()
	case AddressHeatingValve:
		return d.IsFCU4()
	case AddressCoolingValve:
		return d.IsFCU4()
	case AddressHeatingOutput:
		return d.IsHA() || d.IsFCU4()
	case AddressWindowsOpen:
		return d.IsHA() || d.IsFCU4()
	case AddressHoldingFunction, AddressHoldingTimeHi, AddressHoldingTimeLow, AddressHoldingTemperature:
		return d.IsHA()
	case AddressSystemMode:
		return d.IsFCU4()
	case AddressFanSpeedMode, AddressFanSpeed:
		return d.IsFCU4()
	case AddressAway, AddressAwayTemperature:
		return d.IsHA()
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

func (d Device) IsSupportedWindowsOpen() bool {
	return d.IsSupported(AddressWindowsOpen)
}

func (d Device) IsSupportedHeatingOutput() bool {
	return d.IsSupported(AddressHeatingOutput)
}

func (d Device) IsSupportedHoldingFunction() bool {
	return d.IsSupported(AddressHoldingFunction)
}

func (d Device) IsSupportedFloorOverheat() bool {
	return d.IsSupported(AddressFloorOverheat)
}

func (d Device) IsSupportedDeviceType() bool {
	return d.IsSupported(AddressDeviceType)
}

func (d Device) IsSupportedFanSpeedMode() bool {
	return d.IsSupported(AddressFanSpeedMode)
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

func (d Device) IsSupportedHoldingTimeHi() bool {
	return d.IsSupported(AddressHoldingTimeHi)
}

func (d Device) IsSupportedHoldingTimeLow() bool {
	return d.IsSupported(AddressHoldingTimeLow)
}

func (d Device) IsSupportedHoldingTemperature() bool {
	return d.IsSupported(AddressHoldingTemperature)
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

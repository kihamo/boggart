package neptun

const (
	moduleConfigurationFloorWashing uint16 = 1 << iota
	moduleConfigurationFirstGroupAlert
	moduleConfigurationSecondGroupAlert
	moduleConfigurationWirelessSensorLowBattery
	moduleConfigurationWirelessSensorLoss
	moduleConfigurationFirstGroupTapClosing
	moduleConfigurationSecondGroupTapClosing
	moduleConfigurationWirelessSensorPairingMode
	moduleConfigurationFirstGroupTapState
	moduleConfigurationSecondGroupTapState
	moduleConfigurationTwoGroupsMode
	moduleConfigurationTapsClosingOnSensorLoss
	moduleConfigurationKeyboardLock
)

type ModuleConfiguration struct {
	value uint16
}

func (c *ModuleConfiguration) Value() uint16 {
	return c.value
}

func (c *ModuleConfiguration) FloorWashing() bool {
	return c.value&moduleConfigurationFloorWashing != 0
}

func (c *ModuleConfiguration) SetFloorWashing(value bool) {
	if value {
		c.value |= moduleConfigurationFloorWashing
	} else {
		c.value &= ^moduleConfigurationFloorWashing
	}
}

func (c *ModuleConfiguration) FirstGroupAlert() bool {
	return c.value&moduleConfigurationFirstGroupAlert != 0
}

func (c *ModuleConfiguration) SecondGroupAlert() bool {
	return c.value&moduleConfigurationSecondGroupAlert != 0
}

func (c *ModuleConfiguration) WirelessSensorLowBattery() bool {
	return c.value&moduleConfigurationWirelessSensorLowBattery != 0
}

func (c *ModuleConfiguration) WirelessSensorLoss() bool {
	return c.value&moduleConfigurationWirelessSensorLoss != 0
}

func (c *ModuleConfiguration) FirstGroupTapClosing() bool {
	return c.value&moduleConfigurationFirstGroupTapClosing != 0
}

func (c *ModuleConfiguration) SecondGroupTapClosing() bool {
	return c.value&moduleConfigurationSecondGroupTapClosing != 0
}

func (c *ModuleConfiguration) WirelessSensorPairingMode() bool {
	return c.value&moduleConfigurationWirelessSensorPairingMode != 0
}

func (c *ModuleConfiguration) SetWirelessSensorPairingMode(value bool) {
	if value {
		c.value |= moduleConfigurationWirelessSensorPairingMode
	} else {
		c.value &= ^moduleConfigurationWirelessSensorPairingMode
	}
}

func (c *ModuleConfiguration) FirstGroupTapState() bool {
	return c.value&moduleConfigurationFirstGroupTapState != 0
}

func (c *ModuleConfiguration) SetFirsGroupTapState(value bool) {
	if value {
		c.value |= moduleConfigurationFirstGroupTapState
	} else {
		c.value &= ^moduleConfigurationFirstGroupTapState
	}
}

func (c *ModuleConfiguration) SecondGroupTapState() bool {
	return c.value&moduleConfigurationSecondGroupTapState != 0
}

func (c *ModuleConfiguration) SetSecondGroupTapState(value bool) {
	if value {
		c.value |= moduleConfigurationSecondGroupTapState
	} else {
		c.value &= ^moduleConfigurationSecondGroupTapState
	}
}

func (c *ModuleConfiguration) TwoGroupsMode() bool {
	return c.value&moduleConfigurationTwoGroupsMode != 0
}

func (c *ModuleConfiguration) SetTwoGroupsMode(value bool) {
	if value {
		c.value |= moduleConfigurationTwoGroupsMode
	} else {
		c.value &= ^moduleConfigurationTwoGroupsMode
	}
}

func (c *ModuleConfiguration) TapsClosingOnSensorLoss() bool {
	return c.value&moduleConfigurationTapsClosingOnSensorLoss != 0
}

func (c *ModuleConfiguration) SetTapsClosingOnSensorLoss(value bool) {
	if value {
		c.value |= moduleConfigurationTapsClosingOnSensorLoss
	} else {
		c.value &= ^moduleConfigurationTapsClosingOnSensorLoss
	}
}

func (c *ModuleConfiguration) KeyboardLock() bool {
	return c.value&moduleConfigurationKeyboardLock != 0
}

func (c *ModuleConfiguration) SetKeyboardLock(value bool) {
	if value {
		c.value |= moduleConfigurationKeyboardLock
	} else {
		c.value &= ^moduleConfigurationKeyboardLock
	}
}

const (
	inputLinesConfigurationTapFirstGroup uint8 = 1 + iota
	inputLinesConfigurationTapSecondGroup
	inputLinesConfigurationTapTwoGroup
	inputLinesConfigurationType
)

type InputLinesConfiguration struct {
	value uint8
}

func (i *InputLinesConfiguration) Value() uint8 {
	return i.value
}

func (i *InputLinesConfiguration) Tap() uint8 {
	return i.value &^ 0b1100
}

func (i *InputLinesConfiguration) SetTap(value uint8) {
	i.value = (i.Type() << 2) | value
}

func (i *InputLinesConfiguration) TapFirstGroup() bool {
	return i.value&inputLinesConfigurationTapFirstGroup != 0
}

func (i *InputLinesConfiguration) TapSecondGroup() bool {
	return i.value&inputLinesConfigurationTapSecondGroup != 0
}

func (i *InputLinesConfiguration) TapTwoGroup() bool {
	return i.value&inputLinesConfigurationTapTwoGroup == 0b11
}

func (i *InputLinesConfiguration) Type() uint8 {
	return i.value >> 2
}

func (i *InputLinesConfiguration) SetType(value uint8) {
	i.value = (value << 2) | i.Tap()
}

func (i *InputLinesConfiguration) Button() bool {
	return i.value&inputLinesConfigurationType != 0
}

func (i *InputLinesConfiguration) SetButton() {
	i.value |= inputLinesConfigurationType
}

func (i *InputLinesConfiguration) Sensor() bool {
	return i.value&inputLinesConfigurationType == 0
}

func (i *InputLinesConfiguration) SetSensor() {
	i.value &= ^inputLinesConfigurationType
}

const (
	TapSwitchNoSwitch uint8 = iota
	TapSwitchFirstGroup
	TapSwitchSecondGroup
	TapSwitchTwoGroups
)

type EventsRelayConfiguration struct {
	value uint16
}

func (e *EventsRelayConfiguration) Value() uint16 {
	return e.value
}

func (e *EventsRelayConfiguration) TapSwitchOnClosing() uint8 {
	return uint8(e.value >> 2)
}

func (e *EventsRelayConfiguration) SetTapSwitchOnClosing(value uint8) {
	e.value = uint16((value << 2) | e.TapSwitchOnAlert())
}

func (e *EventsRelayConfiguration) TapSwitchOnClosingNoSwitch() bool {
	return e.TapSwitchOnClosing() == TapSwitchNoSwitch
}

func (e *EventsRelayConfiguration) TapSwitchOnClosingFirstGroup() bool {
	return e.TapSwitchOnClosing() == TapSwitchFirstGroup
}

func (e *EventsRelayConfiguration) TapSwitchOnClosingSecondGroup() bool {
	return e.TapSwitchOnClosing() == TapSwitchSecondGroup
}

func (e *EventsRelayConfiguration) TapSwitchOnClosingTwoGroups() bool {
	return e.TapSwitchOnClosing() == TapSwitchTwoGroups
}

func (e *EventsRelayConfiguration) TapSwitchOnAlert() uint8 {
	return uint8(e.value &^ 0b1100)
}

func (e *EventsRelayConfiguration) SetTapSwitchOnAlert(value uint8) {
	e.value = uint16(e.TapSwitchOnClosing()<<2 | value)
}

func (e *EventsRelayConfiguration) TapSwitchOnAlertNoSwitch() bool {
	return e.TapSwitchOnAlert() == TapSwitchNoSwitch
}

func (e *EventsRelayConfiguration) TapSwitchOnAlertFirstGroup() bool {
	return e.TapSwitchOnAlert() == TapSwitchFirstGroup
}

func (e *EventsRelayConfiguration) TapSwitchOnAlertSecondGroup() bool {
	return e.TapSwitchOnAlert() == TapSwitchSecondGroup
}

func (e *EventsRelayConfiguration) TapSwitchOnAlertTwoGroups() bool {
	return e.TapSwitchOnAlert() == TapSwitchTwoGroups
}

type CounterConfiguration struct {
	value uint16
}

func (c *CounterConfiguration) Value() uint16 {
	return c.value
}

func (c *CounterConfiguration) Enabled() bool {
	return c.value&0b1 != 0
}

func (c *CounterConfiguration) Type() uint16 {
	return c.value & 0b10
}

func (c *CounterConfiguration) Error() uint8 {
	return uint8((c.value &^ 0b111111110011) >> 2)
}

func (c *CounterConfiguration) Step() uint8 {
	return uint8(c.value >> 8)
}

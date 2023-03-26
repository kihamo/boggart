package neptun

const (
	moduleConfigurationFloorWashing uint16 = 1 << iota // w
	moduleConfigurationFirstGroupAlert
	moduleConfigurationSecondGroupAlert
	moduleConfigurationWirelessSensorLowBattery
	moduleConfigurationWirelessSensorLoss
	moduleConfigurationFirstGroupTapClosing
	moduleConfigurationSecondGroupTapClosing
	moduleConfigurationWirelessSensorPairingMode
	moduleConfigurationFirstGroupTapState      // w
	moduleConfigurationSecondGroupTapState     // w
	moduleConfigurationTwoGroupsMode           // w
	moduleConfigurationTapsClosingOnSensorLoss // w
	moduleConfigurationKeyboardLock            // w
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

func (c *ModuleConfiguration) FirsGroupTapState() bool {
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

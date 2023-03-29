package neptun

type ModuleConfiguration struct {
	keyboardLock              bool
	tapsClosingOnSensorLoss   bool
	twoGroupsMode             bool
	secondGroupTapState       bool
	firstGroupTapState        bool
	wirelessSensorPairingMode bool
	secondGroupTapClosing     bool
	firstGroupTapClosing      bool
	wirelessSensorLoss        bool
	wirelessSensorLowBattery  bool
	secondGroupAlert          bool
	firstGroupAlert           bool
	floorWashing              bool
}

func NewModuleConfiguration(value uint) *ModuleConfiguration {
	return &ModuleConfiguration{
		keyboardLock:              value&0b1000000000000 != 0,
		tapsClosingOnSensorLoss:   value&0b100000000000 != 0,
		twoGroupsMode:             value&0b10000000000 != 0,
		secondGroupTapState:       value&0b1000000000 != 0,
		firstGroupTapState:        value&0b100000000 != 0,
		wirelessSensorPairingMode: value&0b10000000 != 0,
		secondGroupTapClosing:     value&0b1000000 != 0,
		firstGroupTapClosing:      value&0b100000 != 0,
		wirelessSensorLoss:        value&0b10000 != 0,
		wirelessSensorLowBattery:  value&0b1000 != 0,
		secondGroupAlert:          value&0b100 != 0,
		firstGroupAlert:           value&0b10 != 0,
		floorWashing:              value&0b1 != 0,
	}
}

func (c *ModuleConfiguration) Value() (value uint16) {
	if c.keyboardLock {
		value |= 0b1
	}

	value <<= 1
	if c.tapsClosingOnSensorLoss {
		value |= 0b1
	}

	value <<= 1
	if c.twoGroupsMode {
		value |= 0b1
	}

	value <<= 1
	if c.secondGroupTapState {
		value |= 0b1
	}

	value <<= 1
	if c.firstGroupTapState {
		value |= 0b1
	}

	value <<= 1
	if c.wirelessSensorPairingMode {
		value |= 0b1
	}

	value <<= 1
	if c.secondGroupTapClosing {
		value |= 0b1
	}

	value <<= 1
	if c.firstGroupTapClosing {
		value |= 0b1
	}

	value <<= 1
	if c.wirelessSensorLoss {
		value |= 0b1
	}

	value <<= 1
	if c.wirelessSensorLowBattery {
		value |= 0b1
	}

	value <<= 1
	if c.secondGroupAlert {
		value |= 0b1
	}

	value <<= 1
	if c.firstGroupAlert {
		value |= 0b1
	}

	value <<= 1
	if c.floorWashing {
		value |= 0b1
	}

	return value
}

func (c *ModuleConfiguration) FloorWashing() bool {
	return c.floorWashing
}

func (c *ModuleConfiguration) SetFloorWashing(value bool) {
	c.floorWashing = value
}

func (c *ModuleConfiguration) FirstGroupAlert() bool {
	return c.firstGroupAlert
}

func (c *ModuleConfiguration) SecondGroupAlert() bool {
	return c.secondGroupAlert
}

func (c *ModuleConfiguration) WirelessSensorLowBattery() bool {
	return c.wirelessSensorLowBattery
}

func (c *ModuleConfiguration) WirelessSensorLoss() bool {
	return c.wirelessSensorLoss
}

func (c *ModuleConfiguration) FirstGroupTapClosing() bool {
	return c.firstGroupTapClosing
}

func (c *ModuleConfiguration) SecondGroupTapClosing() bool {
	return c.secondGroupTapClosing
}

func (c *ModuleConfiguration) WirelessSensorPairingMode() bool {
	return c.wirelessSensorPairingMode
}

func (c *ModuleConfiguration) SetWirelessSensorPairingMode(value bool) {
	c.wirelessSensorPairingMode = value
}

func (c *ModuleConfiguration) FirstGroupTapState() bool {
	return c.firstGroupTapState
}

func (c *ModuleConfiguration) SetFirsGroupTapState(value bool) {
	c.firstGroupTapState = value
}

func (c *ModuleConfiguration) SecondGroupTapState() bool {
	return c.secondGroupTapState
}

func (c *ModuleConfiguration) SetSecondGroupTapState(value bool) {
	c.secondGroupTapState = value
}

func (c *ModuleConfiguration) TwoGroupsMode() bool {
	return c.twoGroupsMode
}

func (c *ModuleConfiguration) SetTwoGroupsMode(value bool) {
	c.twoGroupsMode = value
}

func (c *ModuleConfiguration) TapsClosingOnSensorLoss() bool {
	return c.tapsClosingOnSensorLoss
}

func (c *ModuleConfiguration) SetTapsClosingOnSensorLoss(value bool) {
	c.tapsClosingOnSensorLoss = value
}

func (c *ModuleConfiguration) KeyboardLock() bool {
	return c.keyboardLock
}

func (c *ModuleConfiguration) SetKeyboardLock(value bool) {
	c.keyboardLock = value
}

const (
	InputLineTapFirstGroup uint = 1 + iota
	InputLineTapSecondGroup
	InputLineTapTwoGroups
)

const (
	InputLineTypeSensor uint = iota
	InputLineTypeButton
)

type InputLinesConfiguration struct {
	typ uint
	tap uint
}

func NewInputLinesConfiguration(value uint) *InputLinesConfiguration {
	return &InputLinesConfiguration{
		typ: value >> 2,
		tap: value & 0b11,
	}
}

func (i *InputLinesConfiguration) Value() (value uint8) {
	value = uint8(i.typ)

	value <<= 2
	value |= uint8(i.tap)

	return value
}

func (i *InputLinesConfiguration) Tap() uint {
	return i.tap
}

func (i *InputLinesConfiguration) SetTap(value uint) {
	i.tap = value
}

func (i *InputLinesConfiguration) TapFirstGroup() bool {
	return i.tap == InputLineTapFirstGroup
}

func (i *InputLinesConfiguration) TapSecondGroup() bool {
	return i.tap == InputLineTapSecondGroup
}

func (i *InputLinesConfiguration) TapTwoGroup() bool {
	return i.tap == InputLineTapTwoGroups
}

func (i *InputLinesConfiguration) Type() uint {
	return i.typ
}

func (i *InputLinesConfiguration) SetType(value uint) {
	i.typ = value
}

func (i *InputLinesConfiguration) Button() bool {
	return i.typ == InputLineTypeButton
}

func (i *InputLinesConfiguration) SetButton() {
	i.SetType(InputLineTypeButton)
}

func (i *InputLinesConfiguration) Sensor() bool {
	return i.typ == InputLineTypeSensor
}

func (i *InputLinesConfiguration) SetSensor() {
	i.SetType(InputLineTypeSensor)
}

const (
	TapSwitchNoSwitch uint = iota
	TapSwitchFirstGroup
	TapSwitchSecondGroup
	TapSwitchTwoGroups
)

type EventsRelayConfiguration struct {
	closing uint
	alert   uint
}

func NewEventsRelayConfiguration(value uint) *EventsRelayConfiguration {
	return &EventsRelayConfiguration{
		closing: value >> 2,
		alert:   value & 0b11,
	}
}

func (e *EventsRelayConfiguration) Value() (value uint8) {
	value = uint8(e.closing)

	value <<= 2
	value |= uint8(e.alert)

	return value
}

func (e *EventsRelayConfiguration) TapSwitchOnClosing() uint {
	return e.closing
}

func (e *EventsRelayConfiguration) SetTapSwitchOnClosing(value uint) {
	e.closing = value
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

func (e *EventsRelayConfiguration) TapSwitchOnAlert() uint {
	return e.alert
}

func (e *EventsRelayConfiguration) SetTapSwitchOnAlert(value uint) {
	e.alert = value
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

type WirelessSensorStatus struct {
	alert        bool
	lowBattery   bool
	missed       bool
	link         uint
	batteryLevel uint
}

func NewWirelessSensorStatus(value uint) *WirelessSensorStatus {
	return &WirelessSensorStatus{
		alert:        value&0b1 != 0,
		lowBattery:   value&0b10 != 0,
		missed:       value&0b100 != 0,
		link:         value & 0b111000,
		batteryLevel: value & 0b1111111100000000,
	}
}

func (s *WirelessSensorStatus) Alert() bool {
	return s.alert
}

func (s *WirelessSensorStatus) LowBattery() bool {
	return s.lowBattery
}

func (s *WirelessSensorStatus) Missed() bool {
	return s.missed
}

func (s *WirelessSensorStatus) Link() uint {
	return s.link
}

func (s *WirelessSensorStatus) BatteryLevel() uint {
	return s.batteryLevel
}

const (
	CounterTypeBasic uint = iota
	CounterTypeNamur
)

const (
	CounterErrorNo uint = iota
	CounterErrorShortCircuit
	CounterErrorLineBreak
)

type CounterConfiguration struct {
	state bool
	typ   uint
	error uint
	step  uint
}

func NewCounterConfiguration(value uint) *CounterConfiguration {
	cfg := &CounterConfiguration{
		state: value&0b1 != 0,
		typ:   (value >> 1) & 0b1,
		error: (value >> 2) & 0b111111,
	}

	step := value >> 8
	switch step {
	case 1, 10, 100:
		cfg.step = step
	default:
		cfg.step = 10
	}

	return cfg
}

func (c *CounterConfiguration) Value() (value uint16) {
	value = uint16(c.step)

	value <<= 6
	value |= uint16(c.error)

	value <<= 1
	if c.typ == CounterTypeNamur {
		value |= 0b1
	}

	value <<= 1
	if c.state {
		value |= 0b1
	}

	return value
}

func (c *CounterConfiguration) State() bool {
	return c.state
}

func (c *CounterConfiguration) Enabled() bool {
	return c.state
}

func (c *CounterConfiguration) Disabled() bool {
	return c.state == false
}

func (c *CounterConfiguration) SetState(value bool) {
	c.state = value
}

func (c *CounterConfiguration) Disable() {
	c.SetState(false)
}

func (c *CounterConfiguration) Enable() {
	c.SetState(true)
}

func (c *CounterConfiguration) Type() uint {
	return c.typ
}

func (c *CounterConfiguration) SetType(value uint) {
	c.typ = value
}

func (c *CounterConfiguration) SetTypeBasic() {
	c.SetType(CounterTypeBasic)
}

func (c *CounterConfiguration) SetTypeNamur() {
	c.SetType(CounterTypeNamur)
}

func (c *CounterConfiguration) Error() uint {
	return c.error
}

func (c *CounterConfiguration) Step() uint {
	return c.step
}

func (c *CounterConfiguration) SetStep(value uint) {
	switch value {
	case 1, 10, 100:
		c.step = value
	}
}

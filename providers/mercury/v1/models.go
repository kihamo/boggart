package v1

type changed struct {
	changed bool
}

func (c *changed) IsChanged() bool {
	return c.changed
}

func (c *changed) change() {
	c.changed = true
}

type DisplayMode struct {
	changed
	bit uint8
}

func NewDisplayMode(bit uint8) *DisplayMode {
	return &DisplayMode{
		changed: changed{},
		bit:     bit,
	}
}

func (m *DisplayMode) Bit() uint8 {
	return m.bit
}

func (m *DisplayMode) IsTariff1() bool {
	return m.bit&displayModeTariff1 != 0
}

func (m *DisplayMode) SetTariff1(value bool) {
	if value != m.IsTariff1() {
		m.change()
	}

	if value {
		m.bit |= displayModeTariff1
	} else {
		m.bit &= ^displayModeTariff1
	}
}

func (m *DisplayMode) IsTariff2() bool {
	return m.bit&displayModeTariff2 != 0
}

func (m *DisplayMode) SetTariff2(value bool) {
	if value != m.IsTariff2() {
		m.change()
	}

	if value {
		m.bit |= displayModeTariff2
	} else {
		m.bit &= ^displayModeTariff2
	}
}

func (m *DisplayMode) IsTariff3() bool {
	return m.bit&displayModeTariff3 != 0
}

func (m *DisplayMode) SetTariff3(value bool) {
	if value != m.IsTariff3() {
		m.change()
	}

	if value {
		m.bit |= displayModeTariff3
	} else {
		m.bit &= ^displayModeTariff3
	}
}

func (m *DisplayMode) IsTariff4() bool {
	return m.bit&displayModeTariff4 != 0
}

func (m *DisplayMode) SetTariff4(value bool) {
	if value != m.IsTariff4() {
		m.change()
	}

	if value {
		m.bit |= displayModeTariff4
	} else {
		m.bit &= ^displayModeTariff4
	}
}

func (m *DisplayMode) IsAmount() bool {
	return m.bit&displayModeAmount != 0
}

func (m *DisplayMode) SetAmount(value bool) {
	if value != m.IsAmount() {
		m.change()
	}

	if value {
		m.bit |= displayModeAmount
	} else {
		m.bit &= ^displayModeAmount
	}
}

func (m *DisplayMode) IsPower() bool {
	return m.bit&displayModePower != 0
}

func (m *DisplayMode) SetPower(value bool) {
	if value != m.IsPower() {
		m.change()
	}

	if value {
		m.bit |= displayModePower
	} else {
		m.bit &= ^displayModePower
	}
}

func (m *DisplayMode) IsTime() bool {
	return m.bit&displayModeTime != 0
}

func (m *DisplayMode) SetTime(value bool) {
	if value != m.IsTime() {
		m.change()
	}

	if value {
		m.bit |= displayModeTime
	} else {
		m.bit &= ^displayModeTime
	}
}

func (m *DisplayMode) IsDate() bool {
	return m.bit&displayModeDate != 0
}

func (m *DisplayMode) SetDate(value bool) {
	if value != m.IsDate() {
		m.change()
	}

	if value {
		m.bit |= displayModeDate
	} else {
		m.bit &= ^displayModeDate
	}
}

type DisplayModeExt struct {
	bit uint8
}

func NewDisplayModeExt(bit uint8) *DisplayModeExt {
	return &DisplayModeExt{
		bit: bit,
	}
}

func (m *DisplayModeExt) Bit() uint8 {
	return m.bit
}

func (m *DisplayModeExt) IsTariffSchedule() bool {
	return m.bit&displayModeTariffSchedule != 0
}

func (m *DisplayModeExt) IsUIF() bool {
	return m.bit&displayModeUIF != 0
}

func (m *DisplayModeExt) IsReactiveEnergy() bool {
	return m.bit&displayModeReactiveEnergy != 0
}

func (m *DisplayModeExt) IsMaximumResets() bool {
	return m.bit&displayModeMaximumResets != 0
}

func (m *DisplayModeExt) IsWorkingTime() bool {
	return m.bit&displayModeWorkingTime != 0
}

func (m *DisplayModeExt) IsBatteryLifetime() bool {
	return m.bit&displayModeBatteryLifetime != 0
}

func (m *DisplayModeExt) IsPowerLimit() bool {
	return m.bit&displayModePowerLimit != 0
}

func (m *DisplayModeExt) IsEnergyLimit() bool {
	return m.bit&displayModeEnergyLimit != 0
}

type TariffValues struct {
	changed
	t1 uint64
	t2 uint64
	t3 uint64
	t4 uint64
}

func NewTariffValues(t1, t2, t3, t4 uint64) *TariffValues {
	return &TariffValues{
		changed: changed{},
		t1:      t1,
		t2:      t2,
		t3:      t3,
		t4:      t4,
	}
}

func (v *TariffValues) Tariff1() uint64 {
	return v.t1
}

func (v *TariffValues) SetTariff1(value uint64) {
	if value != v.Tariff1() {
		v.change()
	}

	v.t1 = value
}

func (v *TariffValues) Tariff2() uint64 {
	return v.t2
}

func (v *TariffValues) SetTariff2(value uint64) {
	if value != v.Tariff2() {
		v.change()
	}

	v.t2 = value
}

func (v *TariffValues) Tariff3() uint64 {
	return v.t3
}

func (v *TariffValues) SetTariff3(value uint64) {
	if value != v.Tariff3() {
		v.change()
	}

	v.t3 = value
}

func (v *TariffValues) Tariff4() uint64 {
	return v.t4
}

func (v *TariffValues) SetTariff4(value uint64) {
	if value != v.Tariff4() {
		v.change()
	}

	v.t4 = value
}

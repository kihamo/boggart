package scale

import (
	"time"
)

const (
	UnitKG  Unit = 0x02
	UnitKG2 Unit = 0x22
	UnitLBS Unit = 0x03
	UnitJIN Unit = 0x12
)

type Unit byte

func (u Unit) String() string {
	switch u {
	case UnitKG, UnitKG2:
		return "kg"

	case UnitLBS:
		return "lbs"

	case UnitJIN:
		return "jin"

	default:
		return "unknown"
	}
}

type Metrics struct {
	datetime  time.Time
	unit      Unit
	weight    float64
	impedance int64
}

func (m *Metrics) Datetime() time.Time {
	return m.datetime
}

func (m *Metrics) Unit() Unit {
	return m.unit
}

func (m *Metrics) Weight() float64 {
	return m.weight
}

func (m *Metrics) Impedance() int64 {
	return m.impedance
}

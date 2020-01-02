package scale

import (
	"time"
)

type Measure struct {
	datetime  time.Time
	unit      Unit
	weight    float64
	impedance uint64
}

func NewMeasure(datetime time.Time, unit Unit, weight float64, impedance uint64) *Measure {
	return &Measure{
		datetime:  datetime,
		unit:      unit,
		weight:    weight,
		impedance: impedance,
	}
}

func (m *Measure) Datetime() time.Time {
	return m.datetime
}

func (m *Measure) Unit() Unit {
	return m.unit
}

func (m *Measure) Weight() float64 {
	return m.weight
}

func (m *Measure) Impedance() uint64 {
	return m.impedance
}

func (m *Measure) Metrics(sex sex, height, age uint64) (*Metrics, error) {
	return NewMetrics(sex, age, m.weight, height, m.impedance)
}

func (m *Measure) Scales(sex sex, height, age uint64) (*Scales, error) {
	return NewScales(sex, age, m.weight, height)
}

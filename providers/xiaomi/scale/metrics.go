package scale

import (
	"errors"
	"math"
)

type Metrics struct {
	sex       sex
	age       float64
	weight    float64
	height    float64
	impedance float64

	scales *Scales
}

func NewMetrics(sex sex, age uint64, weight float64, height, impedance uint64) (*Metrics, error) {
	scales, err := NewScales(sex, age, weight, height)
	if err != nil {
		return nil, err
	}

	if impedance == 0 || impedance > 3000 {
		return nil, errors.New("impedance is either too low or too high (limits: = 0 ohm and >3000 ohm)")
	}

	return &Metrics{
		sex:       sex,
		age:       float64(age),
		weight:    weight,
		height:    float64(height),
		impedance: float64(impedance),
		scales:    scales,
	}, nil
}

func (m *Metrics) Scales() *Scales {
	return m.scales
}

func (m *Metrics) LBMCoefficient() (value float64) {
	value = (m.height * 9.058 / 100) * (m.height / 100)
	value += m.weight*0.32 + 12.226
	value -= m.impedance * 0.0068
	value -= m.age * 0.0542

	return value
}

func (m *Metrics) BMR() uint64 {
	var value float64

	if m.sex == SexFemale {
		value = 864.6 + m.weight*10.2036
		value -= m.height * 0.39336
		value -= m.age * 6.204

		if value > 2996 {
			value = 5000
		}
	} else {
		value = 877.8 + m.weight*14.916
		value -= m.height * 0.726
		value -= m.age * 8.976

		if value > 2322 {
			value = 5000
		}
	}

	value = checkValueOverflow(value, 500, 10000)

	return uint64(math.Ceil(value))
}

func (m *Metrics) BMI() (value float64) {
	return checkValueOverflow(m.weight/((m.height/100)*(m.height/100)), 10, 90)
}

func (m *Metrics) FatPercentage() (value float64) {
	con := 0.8

	if m.sex == SexFemale {
		if m.age <= 49 {
			con = 9.25
		} else {
			con = 7.25
		}
	}

	lbm := m.LBMCoefficient()

	coefficient := 1.0

	switch v := m.weight; {
	case m.sex == SexMale && v < 61:
		coefficient = .98
	case m.sex == SexFemale && v > 60:
		coefficient = .96

		if v > 160 {
			coefficient *= 1.03
		}
	case m.sex == SexFemale && v < 50:
		coefficient = 1.02

		if v > 160 {
			coefficient *= 1.03
		}
	}

	value = (1.0 - (((lbm - con) * coefficient) / m.weight)) * 100
	if value > 63 {
		value = 75
	}

	return checkValueOverflow(value, 5, 75)
}

func (m *Metrics) WaterPercentage() (value float64) {
	value = (100 - m.FatPercentage()) * 0.7
	coefficient := 0.98

	if value <= 50 {
		coefficient = 1.02
	}

	if value*coefficient >= 65 {
		value = 75
	}

	return checkValueOverflow(value*coefficient, 35, 75)
}

func (m *Metrics) BoneMass() (value float64) {
	base := 0.18016894

	if m.sex == SexFemale {
		base = 0.245691014
	}

	value = (base - (m.LBMCoefficient() * 0.05158)) * -1

	if value > 2.2 {
		value += 0.1
	} else {
		value -= 0.1
	}

	if m.sex == SexFemale && value > 5.1 {
		value = 8
	} else if m.sex == SexMale && value > 5.2 {
		value = 8
	}

	return checkValueOverflow(value, 0.5, 8)
}

func (m *Metrics) MuscleMass() (value float64) {
	value = m.weight - ((m.FatPercentage() * 0.01) * m.weight) - m.BoneMass()

	if m.sex == SexFemale && value >= 84 {
		value = 120
	} else if m.sex == SexMale && value >= 93.5 {
		value = 120
	}

	return checkValueOverflow(value, 10, 120)
}

func (m *Metrics) VisceralFat() (value float64) {
	if m.sex == SexFemale {
		if m.weight > (13-(m.height*0.5))*-1 {
			subsubcalc := ((m.height * 1.45) + (m.height*0.1158)*m.height) - 120
			subcalc := m.weight * 500 / subsubcalc

			value = (subcalc - 6) + (m.age * 0.07)
		} else {
			subcalc := 0.691 + (m.height * -0.0024) + (m.height * -0.0024)

			value = (((m.height * 0.027) - (subcalc * m.weight)) * -1) + (m.age * 0.07) - m.age
		}
	} else {
		if m.height < m.weight*1.6 {
			subcalc := ((m.height * 0.4) - (m.height * (m.height * 0.0826))) * -1

			value = ((m.weight * 305) / (subcalc + 48)) - 2.9 + (m.age * 0.15)
		} else {
			subcalc := 0.765 + m.height*-0.0015

			value = (((m.height * 0.143) - (m.weight * subcalc)) * -1) + (m.age * 0.15) - 5.0
		}
	}

	return checkValueOverflow(value, 1, 50)
}

func (m *Metrics) IdealWeight() float64 {
	return checkValueOverflow((22*m.height)*m.height/10000, 5.5, 198)
}

func (m *Metrics) FatMassToIdeal() (value float64) {
	value = (m.weight * (m.FatPercentage() / 100)) - (m.weight * (m.scales.FatPercentage()[2] / 100))

	if value < 0 {
		value *= -1
	}

	return value
}

func (m *Metrics) ProteinPercentage() (value float64) {
	value = 100 - (math.Floor(m.FatPercentage()*100) / 100)
	value -= math.Floor(m.WaterPercentage()*100) / 100
	value -= math.Floor((m.BoneMass()/m.weight*100)*100) / 100

	return checkValueOverflow(value, 5, 32)
}

func (m *Metrics) MetabolicAge() (value float64) {
	if m.sex == SexFemale {
		value = (m.height * -1.1165) + (m.weight * 1.5784) + (m.age * 0.4615) + (m.impedance * 0.0415) + 83.2548
	} else {
		value = (m.height * -0.7471) + (m.weight * 0.9161) + (m.age * 0.4184) + (m.impedance * 0.0517) + 54.2267
	}

	return checkValueOverflow(value, 15, 80)
}

func (m *Metrics) BodyType() (value uint64) {
	switch v := m.FatPercentage(); {
	case v > m.scales.FatPercentage()[2]:
		value = 0
	case v < m.scales.FatPercentage()[1]:
		value = 2
	default:
		value = 1
	}

	switch v := m.MuscleMass(); {
	case v > m.scales.MuscleMass()[1]:
		value = 2 + (value * 3)
	case v < m.scales.MuscleMass()[0]:
		value *= 3
	default:
		value = 1 + (value * 3)
	}

	return value
}

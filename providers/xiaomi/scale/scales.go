package scale

import (
	"errors"
	"math"
)

var (
	scalesBMR = map[sex][]struct {
		minAge uint64
		value  float64
	}{
		SexMale: {
			{minAge: 30, value: 21.6},
			{minAge: 50, value: 20.07},
			{minAge: 100, value: 19.35},
		},
		SexFemale: {
			{minAge: 30, value: 21.24},
			{minAge: 50, value: 19.53},
			{minAge: 100, value: 18.63},
		},
	}

	scalesBMI = []float64{18.5, 25.0, 28.0, 32.0}

	scalesFatPercentage = map[sex][]struct {
		minAge uint64
		maxAge uint64
		values []float64
	}{
		SexMale: {
			{minAge: 0, maxAge: 17, values: []float64{7.0, 16.0, 25.0, 30.0}},
			{minAge: 18, maxAge: 39, values: []float64{11.0, 17.0, 22.0, 27.0}},
			{minAge: 40, maxAge: 59, values: []float64{12.0, 18.0, 23.0, 28.0}},
			{minAge: 60, maxAge: 100, values: []float64{14.0, 20.0, 25.0, 30.0}},
		},
		SexFemale: {
			{minAge: 0, maxAge: 11, values: []float64{12.0, 21.0, 30.0, 34.0}},
			{minAge: 12, maxAge: 13, values: []float64{15.0, 24.0, 33.0, 37.0}},
			{minAge: 14, maxAge: 15, values: []float64{18.0, 27.0, 36.0, 40.0}},
			{minAge: 16, maxAge: 17, values: []float64{20.0, 28.0, 37.0, 41.0}},
			{minAge: 18, maxAge: 39, values: []float64{21.0, 28.0, 35.0, 40.0}},
			{minAge: 40, maxAge: 59, values: []float64{22.0, 29.0, 36.0, 41.0}},
			{minAge: 60, maxAge: 100, values: []float64{23.0, 30.0, 37.0, 42.0}},
		},
	}

	scalesWaterPercentage = map[sex][]float64{
		SexMale:   {55.0, 65.1},
		SexFemale: {45.0, 60.1},
	}

	scalesBoneMass = map[sex][]struct {
		minWeight float64
		values    []float64
	}{
		SexMale: {
			{minWeight: 75, values: []float64{2.0, 4.2}},
			{minWeight: 60, values: []float64{1.9, 4.1}},
			{minWeight: 0, values: []float64{1.6, 3.9}},
		},
		SexFemale: {
			{minWeight: 60, values: []float64{1.8, 3.9}},
			{minWeight: 45, values: []float64{1.5, 3.8}},
			{minWeight: 0, values: []float64{1.3, 3.6}},
		},
	}

	scalesMuscleMass = map[sex][]struct {
		minHeight float64
		values    []float64
	}{
		SexMale: {
			{minHeight: 170, values: []float64{49.4, 59.5}},
			{minHeight: 160, values: []float64{44.0, 52.5}},
			{minHeight: 0, values: []float64{38.5, 46.6}},
		},
		SexFemale: {
			{minHeight: 160, values: []float64{36.5, 42.6}},
			{minHeight: 150, values: []float64{32.9, 37.6}},
			{minHeight: 0, values: []float64{29.1, 34.8}},
		},
	}

	scalesVisceralFat = []float64{10.0, 15.0}

	scalesProteinPercentage = []float64{16, 20}

	scalesBodyType = map[uint64]string{
		0: "obese",
		1: "overweight",
		2: "thick-set",
		3: "lack-exerscise",
		4: "balanced",
		5: "balanced-muscular",
		6: "skinny",
		7: "balanced-skinny",
		8: "skinny-muscular",
	}
)

type Scales struct {
	sex    sex
	age    uint64
	weight float64
	height float64
}

func NewScales(sex sex, age uint64, weight float64, height uint64) (*Scales, error) {
	if height == 0 || height > 220 {
		return nil, errors.New("height is either too low or too high (limits: <0cm and >220cm)")
	}

	if weight < 10 || weight > 200 {
		return nil, errors.New("weight is either too low or too high (limits: <10kg and >200kg)")
	}

	if age == 0 || age > 99 {
		return nil, errors.New("age is either too low or too high (limits: <= 0 years and > 99 years)")
	}

	return &Scales{
		sex:    sex,
		age:    age,
		weight: weight,
		height: float64(height),
	}, nil
}

func (s *Scales) BMR() uint64 {
	var value float64

	for _, row := range scalesBMR[s.sex] {
		if s.age < row.minAge {
			value = s.weight * row.value
			break
		}
	}

	return uint64(math.Round(value))
}

func (s *Scales) BMI() []float64 {
	return scalesBMI
}

func (s *Scales) FatPercentage() []float64 {
	for _, row := range scalesFatPercentage[s.sex] {
		if s.age >= row.minAge && s.age <= row.maxAge {
			return row.values
		}
	}

	return nil
}

func (s *Scales) WaterPercentage() []float64 {
	return scalesWaterPercentage[s.sex]
}

func (s *Scales) BoneMass() []float64 {
	for _, row := range scalesBoneMass[s.sex] {
		if s.weight >= row.minWeight {
			return row.values
		}
	}

	return nil
}

func (s *Scales) MuscleMass() []float64 {
	for _, row := range scalesMuscleMass[s.sex] {
		if s.height >= row.minHeight {
			return row.values
		}
	}

	return nil
}

func (s *Scales) VisceralFat() []float64 {
	return scalesVisceralFat
}

func (s *Scales) IdealWeight() (value []float64) {
	for _, bmiScale := range s.BMI() {
		value = append(value, (bmiScale*s.height)*s.height/10000)
	}

	return value
}

func (s *Scales) ProteinPercentage() []float64 {
	return scalesProteinPercentage
}

func (s *Scales) BodyType() map[uint64]string {
	return scalesBodyType
}

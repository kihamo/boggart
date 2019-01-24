package atomic

import (
	"math"
)

const (
	float32Null = math.MaxFloat64
)

type Float32Null struct {
	Float64
}

func NewFloat32Null() *Float32Null {
	v := &Float32Null{}
	v.Float64.Set(float32Null)
	return v
}

func NewFloat32NullDefault(value float32) *Float32Null {
	v := &Float32Null{}
	v.Set(value)
	return v
}

func (v *Float32Null) Set(value float32) bool {
	return v.Float64.Set(float64(value))
}

func (v *Float32Null) Load() float32 {
	value := v.Float64.Load()
	if value == float32Null {
		return 0
	}

	return float32(value)
}

func (v *Float32Null) IsNil() bool {
	return v.Float64.Load() == float32Null
}

func (v *Float32Null) Nil() bool {
	return v.Float64.Set(float32Null)
}

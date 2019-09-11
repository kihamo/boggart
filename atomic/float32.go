package atomic

import (
	"math"
	"strconv"
)

type Float32 struct {
	Uint32
}

func NewFloat32() *Float32 {
	return &Float32{}
}

func NewFloat32Default(value float32) *Float32 {
	v := &Float32{}
	v.Uint32.v = math.Float32bits(value)

	return v
}

func (v *Float32) Set(value float32) bool {
	return v.Uint32.Set(math.Float32bits(value))
}

func (v *Float32) Load() float32 {
	return math.Float32frombits(v.Uint32.Load())
}

func (v *Float32) String() string {
	return strconv.FormatFloat(float64(v.Load()), 'f', -1, 32)
}

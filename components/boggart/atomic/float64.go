package atomic

import (
	"math"
)

type Float64 struct {
	Uint64
}

func NewFloat64() *Float64 {
	return &Float64{}
}

func (v *Float64) Set(value float64) bool {
	return v.Uint64.Set(math.Float64bits(value))
}

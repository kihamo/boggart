package mc6

import (
	"time"
)

type value uint16

func (v value) Temperature() float64 {
	return float64(v) / 10
}

func (v value) Duration() time.Duration {
	return time.Duration(v) * time.Minute
}

func (v value) Bool() bool {
	return v == 1
}

func (v value) Uint() uint {
	return uint(value)
}

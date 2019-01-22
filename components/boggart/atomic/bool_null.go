package atomic

import (
	a "sync/atomic"
)

const (
	boolNull = uint32(iota)
	boolTrue
	boolFalse
)

type BoolNull struct {
	v uint32
}

func NewBoolNull() *BoolNull {
	return &BoolNull{
		v: boolNull,
	}
}

func NewBoolNullDefault(value bool) *BoolNull {
	v := NewBoolNull()

	if value {
		v.v = boolTrue
	} else {
		v.v = boolFalse
	}

	return v
}

func (v *BoolNull) Set(value bool) bool {
	var current uint32

	if value {
		current = boolTrue
	} else {
		current = boolFalse
	}

	old := a.SwapUint32(&v.v, current)
	return old != current
}

func (v *BoolNull) Load() bool {
	return a.LoadUint32(&v.v) == boolTrue
}

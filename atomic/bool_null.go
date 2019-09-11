package atomic

import (
	"strconv"
	a "sync/atomic"
)

const (
	boolNull = uint32(iota)
	boolTrue
	boolFalse
)

type BoolNull struct {
	Uint32
}

func NewBoolNull() *BoolNull {
	v := &BoolNull{}
	v.Uint32.Set(boolNull)
	return v
}

func NewBoolNullDefault(value bool) *BoolNull {
	var current uint32

	if value {
		current = boolTrue
	} else {
		current = boolFalse
	}

	v := &BoolNull{}
	v.Uint32.Set(current)
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

func (v *BoolNull) True() bool {
	return v.Set(true)
}

func (v *BoolNull) False() bool {
	return v.Set(false)
}

func (v *BoolNull) Load() bool {
	return a.LoadUint32(&v.v) == boolTrue
}

func (v *BoolNull) IsTrue() bool {
	return v.Load()
}

func (v *BoolNull) IsFalse() bool {
	return !v.Load()
}

func (v *BoolNull) IsNil() bool {
	return v.Uint32.Load() == boolNull
}

func (v *BoolNull) Nil() bool {
	return v.Uint32.Set(boolNull)
}

func (v *BoolNull) String() string {
	value := v.Uint32.Load()

	if value == boolNull {
		return nilString
	}

	return strconv.FormatBool(value == boolTrue)
}

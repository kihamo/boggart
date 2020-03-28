package atomic

import (
	"strconv"
	a "sync/atomic"
)

type Bool struct {
	v uint32
}

func NewBool() *Bool {
	return &Bool{}
}

func NewBoolDefault(value bool) *Bool {
	v := &Bool{}

	if value {
		v.v = 1
	}

	return v
}

func (v *Bool) Set(value bool) bool {
	var current uint32
	if value {
		current = 1
	}

	old := a.SwapUint32(&v.v, current)

	return old != current
}

func (v *Bool) True() bool {
	return v.Set(true)
}

func (v *Bool) False() bool {
	return v.Set(false)
}

func (v *Bool) Load() bool {
	return a.LoadUint32(&v.v) == 1
}

func (v *Bool) IsTrue() bool {
	return v.Load()
}

func (v *Bool) IsFalse() bool {
	return !v.Load()
}

func (v *Bool) String() string {
	return strconv.FormatBool(v.Load())
}

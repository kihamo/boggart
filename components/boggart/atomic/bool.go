package atomic

import (
	a "sync/atomic"
)

type Bool struct {
	v uint32
}

func NewBool() *Bool {
	return &Bool{}
}

func (v *Bool) Set(value bool) bool {
	var current uint32
	if value {
		current = 1
	}

	old := a.SwapUint32(&v.v, current)
	return old != current
}

func (v *Bool) Load() bool {
	return a.LoadUint32(&v.v) == 1
}

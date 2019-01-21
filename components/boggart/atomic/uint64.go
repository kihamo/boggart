package atomic

import (
	a "sync/atomic"
)

type Uint64 struct {
	v uint64
}

func NewUint64() *Uint64 {
	return &Uint64{}
}

func (v *Uint64) Set(value uint64) bool {
	old := a.SwapUint64(&v.v, value)
	return old != value
}

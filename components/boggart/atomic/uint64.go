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

func (v *Uint64) Load() uint64 {
	return a.LoadUint64(&v.v)
}

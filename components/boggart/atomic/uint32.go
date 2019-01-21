package atomic

import (
	a "sync/atomic"
)

type Uint32 struct {
	v uint32
}

func NewUint32() *Uint32 {
	return &Uint32{}
}

func NewUint32Default(value uint32) *Uint32 {
	return &Uint32{
		v: value,
	}
}

func (v *Uint32) Set(value uint32) bool {
	old := a.SwapUint32(&v.v, value)
	return old != value
}

func (v *Uint32) Load() uint32 {
	return a.LoadUint32(&v.v)
}

package atomic

import (
	a "sync/atomic"
)

type Int64 struct {
	v int64
}

func NewInt64() *Int64 {
	return &Int64{}
}

func (v *Int64) Set(value int64) bool {
	old := a.SwapInt64(&v.v, value)
	return old != value
}

func (v *Int64) Load() int64 {
	return a.LoadInt64(&v.v)
}

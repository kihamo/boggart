package atomic

import (
	"strconv"
	a "sync/atomic"
)

type Uint64 struct {
	v uint64
}

func NewUint64() *Uint64 {
	return &Uint64{}
}

func NewUint64Default(value uint64) *Uint64 {
	return &Uint64{
		v: value,
	}
}

func (v *Uint64) Set(value uint64) bool {
	old := a.SwapUint64(&v.v, value)
	return old != value
}

func (v *Uint64) Load() uint64 {
	return a.LoadUint64(&v.v)
}

func (v *Uint64) String() string {
	return strconv.FormatUint(v.Load(), 10)
}

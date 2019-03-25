package atomic

import (
	"strconv"
	a "sync/atomic"
)

type Int32 struct {
	v int32
}

func NewInt32() *Int32 {
	return &Int32{}
}

func NewInt32Default(value int32) *Int32 {
	return &Int32{
		v: value,
	}
}

func (v *Int32) Set(value int32) bool {
	old := a.SwapInt32(&v.v, value)
	return old != value
}

func (v *Int32) Load() int32 {
	return a.LoadInt32(&v.v)
}

func (v *Int32) String() string {
	return strconv.FormatInt(int64(v.Load()), 10)
}

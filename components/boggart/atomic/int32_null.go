package atomic

import (
	"math"
)

const (
	int32Null = math.MaxInt64
)

type Int32Null struct {
	Int64
}

func NewInt32Null() *Int32Null {
	v := &Int32Null{}
	v.Int64.Set(int32Null)
	return v
}

func NewInt32NullDefault(value int32) *Int32Null {
	v := &Int32Null{}
	v.Set(value)
	return v
}

func (v *Int32Null) Set(value int32) bool {
	return v.Int64.Set(int64(value))
}

func (v *Int32Null) Load() int32 {
	value := v.Int64.Load()
	if value == int32Null {
		return 0
	}

	return int32(value)
}

func (v *Int32Null) IsNil() bool {
	return v.Int64.Load() == int32Null
}

func (v *Int32Null) Nil() bool {
	return v.Int64.Set(int32Null)
}

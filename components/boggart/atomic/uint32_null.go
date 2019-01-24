package atomic

import (
	"math"
)

const (
	uint32Null = math.MaxUint64
)

type Uint32Null struct {
	Uint64
}

func NewUint32Null() *Uint32Null {
	v := &Uint32Null{}
	v.Uint64.Set(uint32Null)
	return v
}

func NewUint32NullDefault(value uint32) *Uint32Null {
	v := &Uint32Null{}
	v.Set(value)
	return v
}

func (v *Uint32Null) Set(value uint32) bool {
	return v.Uint64.Set(uint64(value))
}

func (v *Uint32Null) Load() uint32 {
	value := v.Uint64.Load()
	if value == uint32Null {
		return 0
	}

	return uint32(value)
}

func (v *Uint32Null) IsNil() bool {
	return v.Uint64.Load() == uint32Null
}

func (v *Uint32Null) Nil() bool {
	return v.Uint64.Set(uint32Null)
}

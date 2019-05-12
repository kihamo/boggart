package atomic

import (
	"strconv"
	"sync"
)

type Int64Null struct {
	m sync.RWMutex
	n bool
	v int64
}

func NewInt64Null() *Int64Null {
	return &Int64Null{
		n: true,
		v: 0,
	}
}

func NewInt64NullDefault(value int64) *Int64Null {
	return &Int64Null{
		n: false,
		v: value,
	}
}

func (v *Int64Null) Set(value int64) bool {
	v.m.Lock()

	oldNull := v.n
	oldValue := v.v

	v.n = false
	v.v = value

	v.m.Unlock()

	return oldNull || oldValue != value
}

func (v *Int64Null) Load() int64 {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.v
}

func (v *Int64Null) IsNil() bool {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.n
}

func (v *Int64Null) Nil() bool {
	v.m.Lock()

	old := v.n

	v.n = true
	v.v = 0

	v.m.Unlock()

	return !old
}

func (v *Int64Null) String() string {
	v.m.RLock()
	defer v.m.RUnlock()

	if v.n {
		return nilString
	}

	return strconv.FormatInt(v.v, 10)
}

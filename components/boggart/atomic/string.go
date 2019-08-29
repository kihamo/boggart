package atomic

import (
	"sync"
)

type String struct {
	m sync.RWMutex
	v string
}

func NewString() *String {
	return &String{}
}

func NewStringDefault(value string) *String {
	return &String{
		v: value,
	}
}

func (v *String) Set(value string) bool {
	v.m.Lock()
	old := v.v
	v.v = value
	v.m.Unlock()

	return old != value
}

func (v *String) Load() string {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.v
}

func (v *String) IsEmpty() bool {
	return v.Load() == ""
}

func (v *String) String() string {
	return v.Load()
}

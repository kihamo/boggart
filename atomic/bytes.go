package atomic

import (
	"bytes"
	"sync"
)

type Bytes struct {
	m sync.RWMutex
	v []byte
}

func NewBytes() *Bytes {
	return &Bytes{}
}

func NewBytesDefault(value []byte) *Bytes {
	return &Bytes{
		v: value,
	}
}

func (v *Bytes) Set(value []byte) bool {
	v.m.Lock()
	old := v.v
	v.v = value
	v.m.Unlock()

	return !bytes.Equal(old, value)
}

func (v *Bytes) Load() []byte {
	v.m.RLock()
	defer v.m.RUnlock()

	return v.v
}

func (v *Bytes) IsEmpty() bool {
	return len(v.Load()) == 0
}

func (v *Bytes) String() string {
	return string(v.Load())
}

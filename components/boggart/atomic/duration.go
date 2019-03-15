package atomic

import (
	"time"
)

type Duration struct {
	Int64
}

func NewDuration() *Duration {
	return &Duration{}
}

func NewDurationDefault(value time.Duration) *Duration {
	v := &Duration{}
	v.Int64.v = int64(value)

	return v
}

func (v *Duration) Set(value time.Duration) bool {
	return v.Int64.Set(int64(value))
}

func (v *Duration) Load() time.Duration {
	return time.Duration(v.Int64.Load())
}

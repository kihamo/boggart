package atomic

import (
	"time"
)

type Time struct {
	Int64
}

func NewTime() *Time {
	return &Time{}
}

func NewTimeDefault(value time.Time) *Time {
	v := &Time{}
	v.Int64.v = value.UnixNano()

	return v
}

func (v *Time) Set(value time.Time) bool {
	return v.Int64.Set(value.UnixNano())
}

func (v *Time) Load() time.Time {
	ns := v.Int64.Load()
	return time.Unix(0, ns)
}

func (v *Time) String() string {
	return v.Load().String()
}

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
	v.Set(value)

	return v
}

func (v *Time) Set(value time.Time) bool {
	if value.IsZero() {
		return v.Int64.Set(0)
	}

	return v.Int64.Set(value.UnixNano())
}

func (v *Time) Load() time.Time {
	ns := v.Int64.Load()
	if ns == 0 {
		return time.Time{}
	}

	return time.Unix(0, ns)
}

func (v *Time) String() string {
	return v.Load().String()
}

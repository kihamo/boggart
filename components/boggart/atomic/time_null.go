package atomic

import (
	"time"
)

var (
	timeNull = time.Time{}
)

type TimeNull struct {
	Time
}

func NewTimeNull() *TimeNull {
	v := &TimeNull{}
	v.Set(timeNull)
	return v
}

func NewTimeNullDefault(value time.Time) *TimeNull {
	v := &TimeNull{}
	v.Set(value)
	return v
}

func (v *TimeNull) IsNil() bool {
	return v.Time.Load() == timeNull
}

func (v *TimeNull) Nil() bool {
	return v.Time.Set(timeNull)
}

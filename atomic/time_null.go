package atomic

import (
	"time"
)

type TimeNull struct {
	Time
}

func NewTimeNull() *TimeNull {
	v := &TimeNull{}
	v.Set(time.Time{})
	return v
}

func NewTimeNullDefault(value time.Time) *TimeNull {
	v := &TimeNull{}
	v.Set(value)
	return v
}

func (v *TimeNull) Load() *time.Time {
	value := v.Time.Load()
	if value.IsZero() {
		return nil
	}

	return &value
}

func (v *TimeNull) IsNil() bool {
	return v.Time.Load().IsZero()
}

func (v *TimeNull) Nil() bool {
	return v.Time.Set(time.Time{})
}

func (v *TimeNull) String() string {
	value := v.Load()
	if value == nil {
		return nilString
	}

	return value.String()
}

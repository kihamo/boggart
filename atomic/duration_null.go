package atomic

import (
	"time"
)

type DurationNull struct {
	Int64Null
}

func NewDurationNull() *DurationNull {
	return &DurationNull{
		Int64Null: *NewInt64Null(),
	}
}

func NewDurationNullDefault(value time.Duration) *DurationNull {
	return &DurationNull{
		Int64Null: *NewInt64NullDefault(int64(value)),
	}
}

func (v *DurationNull) Set(value time.Duration) bool {
	return v.Int64Null.Set(int64(value))
}

func (v *DurationNull) Load() *time.Duration {
	v.m.RLock()
	defer v.m.RUnlock()

	if v.n {
		return nil
	}

	value := time.Duration(v.v)
	return &value
}

func (v *DurationNull) String() string {
	value := v.Load()
	if value == nil {
		return nilString
	}

	return value.String()
}

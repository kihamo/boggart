package gpio

import (
	"time"
)

type GPIOPin interface {
	Number() int64
	Status() bool
	ChangedAt() *time.Time
	SetCallbackChange(callback func(bool, time.Time, *time.Time))
}

type PinMode int64

const (
	PIN_IN PinMode = iota
	PIN_OUT
	PIN_PWM
)

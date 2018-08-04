package gpio

import (
	"time"
)

type GPIOPin interface {
	Up() error
	Down() error
	Mode() (PinMode, error)
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

func (m PinMode) String() string {
	switch m {
	case PIN_IN:
		return "in"

	case PIN_OUT:
		return "out"

	case PIN_PWM:
		return "pwm"
	}

	return "unknown"
}

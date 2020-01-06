package rpi

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("raspberry_pi", Type{})
}

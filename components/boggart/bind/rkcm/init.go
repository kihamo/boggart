package rkcm

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("rkcm", Type{})
}

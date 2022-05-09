package apcupsd

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("apcupsd", Type{})
}

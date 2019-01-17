package alsa

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("alsa", Type{})
}

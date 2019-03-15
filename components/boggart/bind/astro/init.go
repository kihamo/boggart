package astro

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/astro/sun"
)

func init() {
	boggart.RegisterBindType("astro:sun", sun.Type{})
}

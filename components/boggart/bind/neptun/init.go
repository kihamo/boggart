package neptun

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("neptun", Type{}, "neptun:smart")
}

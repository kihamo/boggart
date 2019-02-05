package boggart

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType(boggart.ComponentName, Type{})
}

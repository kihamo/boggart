package mikrotik

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("mikrotik", Type{})
}

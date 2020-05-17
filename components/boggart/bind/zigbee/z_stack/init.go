package z_stack

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("zigbee:z-stack:coordinator", Type{})
}

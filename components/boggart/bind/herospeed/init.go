package herospeed

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("herospeed:ipc", Type{})
}

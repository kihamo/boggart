package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("ds18b20:w1", Type{})
}

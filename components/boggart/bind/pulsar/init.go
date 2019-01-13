package pulsar

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("pulsar:heat_meter", Type{})
}

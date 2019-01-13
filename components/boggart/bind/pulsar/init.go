package pulsar

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("pulsar:heat_meter", Type{})
}

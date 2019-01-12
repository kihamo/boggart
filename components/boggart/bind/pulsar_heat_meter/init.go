package pulsar_heat_meter

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("pulsar_heat_meter", Type{})
}

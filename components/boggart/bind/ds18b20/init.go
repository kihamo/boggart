package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("ds18b20:w1", Type{})
}

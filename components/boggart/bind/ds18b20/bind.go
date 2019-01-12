package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	lastValue int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT
}

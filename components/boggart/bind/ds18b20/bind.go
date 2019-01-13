package ds18b20

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	lastValue int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	livenessInterval time.Duration
	livenessTimeout  time.Duration
	updaterInterval  time.Duration
}

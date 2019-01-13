package ds18b20

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	lastValue int64

	boggart.BindBase
	boggart.BindMQTT

	livenessInterval time.Duration
	livenessTimeout  time.Duration
	updaterInterval  time.Duration
}

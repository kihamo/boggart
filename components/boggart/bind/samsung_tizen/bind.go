package samsung_tizen

import (
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
)

type Bind struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	client           *tv.ApiV2
	mac              string
	livenessInterval time.Duration
	livenessTimeout  time.Duration
}

package samsung_tizen

import (
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	mutex sync.RWMutex

	client           *tv.ApiV2
	mac              string
	livenessInterval time.Duration
	livenessTimeout  time.Duration
}

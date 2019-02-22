package ping

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	hostname        string
	retry           int
	timeout         time.Duration
	updaterInterval time.Duration

	online  *atomic.BoolNull
	latency *atomic.Uint32Null
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}

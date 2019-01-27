package network

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type BindPing struct {
	boggart.BindBase
	boggart.BindMQTT

	hostname        string
	retry           int
	timeout         time.Duration
	updaterInterval time.Duration

	online  *atomic.BoolNull
	latency *atomic.Uint32Null
}

func (b *BindPing) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}

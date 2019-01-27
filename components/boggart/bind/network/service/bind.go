package service

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	address         string
	retry           int
	timeout         time.Duration
	updaterInterval time.Duration

	online  *atomic.BoolNull
	latency *atomic.Uint32Null
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}

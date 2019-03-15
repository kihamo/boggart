package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/go-workers/task"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config

	sunrise  *atomic.TimeNull
	sunset   *atomic.TimeNull
	dayLight *atomic.Duration

	taskStateUpdater *task.FunctionTask
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}

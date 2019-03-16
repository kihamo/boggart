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

	riseStart     *atomic.TimeNull
	riseEnd       *atomic.TimeNull
	riseDuration  *atomic.Duration
	setStart      *atomic.TimeNull
	setEnd        *atomic.TimeNull
	setDuration   *atomic.Duration
	nightStart    *atomic.TimeNull
	nightEnd      *atomic.TimeNull
	nightDuration *atomic.Duration
	nadir         *atomic.TimeNull
	solarNoon     *atomic.TimeNull

	taskStateUpdater *task.FunctionTask
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}

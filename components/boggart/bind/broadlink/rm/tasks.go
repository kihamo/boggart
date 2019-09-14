package rm

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/sparrc/go-ping"
)

const (
	overhead = time.Millisecond * 100
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness-" + b.SerialNumber())

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	pinger, err := ping.NewPinger(b.config.Host)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	pinger.SetPrivileged(true)

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, nil
	}

	pinger.Timeout = time.Until(deadline)
	if pinger.Timeout > overhead {
		pinger.Timeout -= overhead
	}

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv != 0 {
		b.UpdateStatus(boggart.BindStatusOnline)
	} else {
		b.UpdateStatus(boggart.BindStatusOffline)
	}

	return nil, nil
}

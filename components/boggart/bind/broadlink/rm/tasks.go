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
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("bind-broadlink:rm-liveness-" + b.SerialNumber())

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	pinger, err := ping.NewPinger(b.ip.IP.String())
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, nil
	}

	pinger.Timeout = deadline.Sub(time.Now())
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

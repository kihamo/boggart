package chromecast

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("liveness-" + b.host.String())

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	err := b.Connect(ctx)
	if err == nil {
		b.UpdateStatus(boggart.BindStatusOnline)
	} else {
		b.UpdateStatus(boggart.BindStatusOffline)

		b.Logger().Error("Liveness checker failed", "error", err.Error())
	}

	return nil, nil
}

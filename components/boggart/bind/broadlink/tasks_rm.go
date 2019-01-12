package broadlink

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *BindRM) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-broadlink-rm-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (b *BindRM) taskLiveness(ctx context.Context) (interface{}, error) {
	// TODO:
	b.UpdateStatus(boggart.DeviceStatusOnline)

	return nil, nil
}

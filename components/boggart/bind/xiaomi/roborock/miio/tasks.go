package miio

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("bind-xiaomi-roborock-" + b.config.Host)

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	sn, err := b.device.SerialNumber(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if b.Status() == boggart.BindStatusOnline {
		return nil, nil
	}

	if b.SerialNumber() == "" {
		b.SetSerialNumber(sn)
	}
	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, nil
}

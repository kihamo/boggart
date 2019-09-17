package esphome

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
	taskLiveness.SetName("liveness-" + b.config.Address)

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	info, err := b.provider.DeviceInfo(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	if sn := b.SerialNumber(); sn == "" {
		b.SetSerialNumber(info.MacAddress)
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	return nil, nil
}

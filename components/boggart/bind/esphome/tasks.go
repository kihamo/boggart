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

	taskUpdater := task.NewFunctionTask(b.taskUpdated)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater-" + b.config.Address)

	return []workers.Task{
		taskLiveness,
		taskUpdater,
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

func (b *Bind) taskUpdated(ctx context.Context) (interface{}, error) {
	sn := b.SerialNumber()
	if sn == "" {
		return nil, nil
	}

	entities, err := b.provider.ListEntities(ctx)
	if err != nil {
		return nil, err
	}

	return nil, b.syncState(ctx, entities...)
}

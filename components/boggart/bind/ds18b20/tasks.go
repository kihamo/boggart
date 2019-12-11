package ds18b20

import (
	"context"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if !b.IsStatusOnline() {
		return nil, nil
	}

	value, err := b.Temperature()
	if err != nil {
		return nil, err
	}

	metricValue.With("serial_number", b.SerialNumber()).Set(value)

	if err := b.MQTTPublishAsync(ctx, b.config.TopicValue, value); err != nil {
		return nil, err
	}

	return nil, nil
}

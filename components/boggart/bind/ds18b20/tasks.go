package ds18b20

import (
	"context"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	value, err := b.Temperature()
	if err != nil {
		return err
	}

	metricValue.With("serial_number", b.SerialNumber()).Set(value)

	if err := b.MQTTPublishAsync(ctx, b.config.TopicValue, value); err != nil {
		return err
	}

	return nil
}

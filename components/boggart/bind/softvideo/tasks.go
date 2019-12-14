package softvideo

import (
	"context"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	value, err := b.provider.Balance(ctx)
	if err != nil {
		return err
	}

	sn := b.SerialNumber()
	metricBalance.With("account", sn).Set(value)

	err = b.MQTTPublishAsync(ctx, b.config.TopicBalance, value)

	return err
}

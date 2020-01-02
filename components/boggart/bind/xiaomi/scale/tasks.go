package scale

import (
	"context"

	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
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
	metrics, err := b.provider.Metrics(ctx)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		if e := b.MQTTPublishAsync(ctx, b.config.TopicWeight, metric.Weight()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicImpedance, metric.Impedance()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

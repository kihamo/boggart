package openweathermap

import (
	"context"

	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	response, err := b.OneCall(ctx, []string{"current", "daily"})
	if err != nil {
		return err
	}

	id := b.Meta().ID()

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCurrent, response.Current.Temp); e != nil {
		err = multierr.Append(err, e)
	}

	for i, day := range response.Daily {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempMin.Format(id, i), day.Temp.Min); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempMax.Format(id, i), day.Temp.Max); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

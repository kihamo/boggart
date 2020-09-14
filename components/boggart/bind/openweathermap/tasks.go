package openweathermap

import (
	"context"
	"strconv"

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

	metricCurrent.With("id", id).Set(response.Current.Temp)
	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCurrent, response.Current.Temp); e != nil {
		err = multierr.Append(err, e)
	}

	for i, day := range response.Daily {
		dayAsString := strconv.Itoa(i)

		metricTempMin.With("id", id).With("day", dayAsString).Set(day.Temp.Min)
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempMin.Format(id, i), day.Temp.Min); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempMax.With("id", id).With("day", dayAsString).Set(day.Temp.Max)
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempMax.Format(id, i), day.Temp.Max); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempDay.With("id", id).With("day", dayAsString).Set(day.Temp.Day)
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempDay.Format(id, i), day.Temp.Day); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempNight.With("id", id).With("day", dayAsString).Set(day.Temp.Night)
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempNight.Format(id, i), day.Temp.Night); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempMorning.With("id", id).With("day", dayAsString).Set(day.Temp.Morn)
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyTempMorning.Format(id, i), day.Temp.Morn); e != nil {
			err = multierr.Append(err, e)
		}

		metricWindSpeed.With("id", id).With("day", dayAsString).Set(day.WindSpeed)
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDailyWindSpeed.Format(id, i), day.WindSpeed); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

package ds18b20

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdater),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	value, err := b.Temperature()
	if err != nil {
		return err
	}

	metricValue.With("serial_number", b.config.Address).Set(value)

	if err := b.MQTT().PublishAsync(ctx, b.config.TopicValue, value); err != nil {
		return err
	}

	return nil
}

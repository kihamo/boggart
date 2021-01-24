package ds18b20

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	values, err := b.Temperatures()
	if err != nil {
		return err
	}

	for sensor, value := range values {
		metricValue.With("serial_number", sensor).Set(value)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicValue.Format(sensor), value); e != nil {
			err = fmt.Errorf("publish value for sensor %s return error: %w", sensor, e)
		}
	}

	return err
}

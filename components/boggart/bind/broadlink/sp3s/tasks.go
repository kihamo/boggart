package sp3s

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
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

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	state, err := b.State()
	if err != nil {
		return err
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicState, state); e != nil {
		err = multierr.Append(err, e)
	}

	if power, e := b.Power(); e == nil {
		metricPower.With("serial_number", b.config.MAC.String()).Set(power)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicPower, power); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}

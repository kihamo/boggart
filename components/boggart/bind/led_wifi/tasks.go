package ledwifi

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
)

const (
	taskIDUpdater = "updater"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName(taskIDUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	state, err := b.bulb.State(ctx)
	if err != nil {
		return err
	}

	id := b.Meta().ID()
	cfg := b.config()

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStatePower.Format(id), state.Power); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateMode.Format(id), state.Mode); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateSpeed.Format(id), state.Speed); e != nil {
		err = multierr.Append(err, e)
	}

	// in HEX
	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateColor.Format(id), state.Color.String()); e != nil {
		err = multierr.Append(err, e)
	}

	// in HSV
	h, s, v := state.Color.HSV()
	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateColorHSV.Format(id), fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) runUpdater(ctx context.Context) error {
	return b.Workers().TaskRunByName(ctx, taskIDUpdater)
}

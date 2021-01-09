package influxdb

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("execute").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskExecuteHandler),
						b.config.ExecuteTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.ExecuteInterval)),
	}
}

func (b *Bind) taskExecuteHandler(ctx context.Context) error {
	api := b.client.QueryAPI(b.config.Organization)

	result, err := api.Query(ctx, b.config.Query)
	if err != nil {
		return err
	}

	if err = result.Err(); err != nil {
		return err
	}

	for result.Next() {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicResult.Format(b.Meta().ID()), result.Record().Value()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

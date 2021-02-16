package influxdb

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	cfg := b.config()

	return []tasks.Task{
		tasks.NewTask().
			WithName("execute").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskExecuteHandler),
						cfg.ExecuteTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.ExecuteInterval)),
	}
}

func (b *Bind) taskExecuteHandler(ctx context.Context) error {
	cfg := b.config()
	api := b.client.QueryAPI(cfg.Organization)

	result, err := api.Query(ctx, cfg.Query)
	if err != nil {
		return err
	}

	if err = result.Err(); err != nil {
		return err
	}

	for result.Next() {
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicResult.Format(b.Meta().ID()), result.Record().Value()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

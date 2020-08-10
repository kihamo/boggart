package influxdb

import (
	"context"

	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskExecute := b.Workers().WrapTaskIsOnline(b.taskExecute)
	taskExecute.SetTimeout(b.config.ExecuteTimeout)
	taskExecute.SetRepeats(-1)
	taskExecute.SetRepeatInterval(b.config.ExecuteInterval)
	taskExecute.SetName("execute")

	return []workers.Task{
		taskExecute,
	}
}

func (b *Bind) taskExecute(ctx context.Context) (err error) {
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

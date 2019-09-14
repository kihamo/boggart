package rkcm

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/rkcm/client/mobile"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater-" + b.config.Login)

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	params := mobile.NewGetDebtParamsWithContext(ctx).
		WithPhone(b.config.Login)

	response, err := b.client.Mobile.GetDebt(params)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	for _, debt := range response.Payload.Data {
		metricBalance.With("ident", debt.Ident).Set(debt.Sum)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicBalance.Format(mqtt.NameReplace(debt.Ident)), debt.Sum); e != nil {
			err = multierr.Append(e, err)
		}
	}

	return nil, err
}

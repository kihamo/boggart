package elektroset

import (
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("bind-elektroset-updater-" + b.config.Login)

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	accounts, err := b.client.Accounts(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	for _, account := range accounts {
		var totalBalance float64

		for _, service := range account.Services {
			totalBalance += service.Balance
			serviceID := strconv.FormatUint(service.ID, 10)

			metricServiceBalance.With("account", account.Number, "service", serviceID).Set(service.Balance)

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicServiceBalance.Format(account.Number, serviceID), service.Balance); e != nil {
				err = multierr.Append(e, err)
			}
		}

		metricBalance.With("account", account.Number).Set(totalBalance)

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBalance.Format(account.Number), totalBalance); e != nil {
			err = multierr.Append(e, err)
		}
	}

	return nil, err
}

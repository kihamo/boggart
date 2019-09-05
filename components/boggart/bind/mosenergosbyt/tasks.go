package mosenergosbyt

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
	taskUpdater.SetName("bind-mosenergosbyt-updater-" + b.config.Login)

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
		if balance, e := b.client.CurrentBalance(ctx, account.VLProvider.IDAbonent); e != nil {
			err = multierr.Append(e, err)
		} else {
			id := strconv.FormatUint(account.VLProvider.IDAbonent, 10)
			metricBalance.With("abonent", id).Set(balance.Balance)

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBalance.Format(id), balance.Balance); e != nil {
				err = multierr.Append(e, err)
			}
		}
	}

	return nil, err
}

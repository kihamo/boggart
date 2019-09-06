package mosenergosbyt

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
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
	accounts, err := b.client.Accounts(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	for _, account := range accounts {
		if balance, e := b.client.CurrentBalance(ctx, account.Provider.IDAbonent); e != nil {
			err = multierr.Append(e, err)
		} else {
			accountID := strconv.FormatUint(account.Provider.IDAbonent, 10)

			metricBalance.With("account", accountID).Set(balance.Balance)

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBalance.Format(accountID), balance.Balance); e != nil {
				err = multierr.Append(e, err)
			}

			for i, service := range balance.Services {
				var serviceID string

				if id, ok := services[strings.ToLower(service.Service)]; ok {
					serviceID = id
				} else {
					serviceID = strconv.FormatInt(int64(i), 10)
				}

				metricServiceBalance.With("account", accountID, "service", serviceID).Set(service.Balance)

				if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicServiceBalance.Format(accountID, serviceID), service.Balance); e != nil {
					err = multierr.Append(e, err)
				}
			}
		}
	}

	return nil, err
}

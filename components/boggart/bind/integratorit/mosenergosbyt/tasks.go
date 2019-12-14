package mosenergosbyt

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	accounts, err := b.client.Accounts(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if balance, e := b.client.CurrentBalance(ctx, account.Provider.IDAbonent); e != nil {
			err = multierr.Append(e, err)
		} else {
			accountID := strconv.FormatUint(account.Provider.IDAbonent, 10)

			metricBalance.With("account", accountID).Set(balance.Balance)

			if e := b.MQTTPublishAsync(ctx, b.config.TopicBalance.Format(accountID), balance.Balance); e != nil {
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

				if e := b.MQTTPublishAsync(ctx, b.config.TopicServiceBalance.Format(accountID, serviceID), service.Balance); e != nil {
					err = multierr.Append(e, err)
				}
			}
		}
	}

	return err
}

package elektroset

import (
	"context"
	"strconv"
	"time"

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

	dateStart := time.Now().Add(-time.Hour * 24 * 31)
	dateEnd := time.Now()

	for _, account := range accounts {
		var totalBalance float64

		for _, service := range account.Services {
			totalBalance += service.Balance
			serviceID := strconv.FormatUint(service.ID, 10)

			metricServiceBalance.With("account", account.Number, "service", serviceID).Set(service.Balance)

			if e := b.MQTTPublishAsync(ctx, b.config.TopicServiceBalance.Format(account.Number, serviceID), service.Balance); e != nil {
				err = multierr.Append(err, e)
			}

			// balance
			if balances, e := b.client.BalanceDetails(ctx, account.Number, service.Provider, dateStart, dateEnd); e == nil {
				for _, balance := range balances {
					if balance.ValueT1 != nil {
						metricServiceBalance.With("account", account.Number, "service", serviceID, "tariff", "1").Set(*balance.ValueT1)

						if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValueT1.Format(account.Number, serviceID), *balance.ValueT1); e != nil {
							err = multierr.Append(err, e)
						}
					}

					if balance.ValueT2 != nil {
						metricServiceBalance.With("account", account.Number, "service", serviceID, "tariff", "2").Set(*balance.ValueT2)

						if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValueT2.Format(account.Number, serviceID), *balance.ValueT2); e != nil {
							err = multierr.Append(err, e)
						}
					}

					if balance.ValueT3 != nil {
						metricServiceBalance.With("account", account.Number, "service", serviceID, "tariff", "3").Set(*balance.ValueT3)

						if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValueT3.Format(account.Number, serviceID), *balance.ValueT3); e != nil {
							err = multierr.Append(err, e)
						}
					}
				}
			} else {
				err = multierr.Append(err, e)
			}
		}

		metricBalance.With("account", account.Number).Set(totalBalance)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicBalance.Format(account.Number), totalBalance); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}

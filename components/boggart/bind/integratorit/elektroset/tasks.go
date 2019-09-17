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

	dateStart := time.Now().Add(-b.config.BalanceDetailsInterval)
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
				var (
					t1, t2, t3 *float64
					t1Date, t2Date, t3Date *time.Time
				)

				for _, balance := range balances {
					if balance.DatetimeEntity == nil {
						continue
					}

					if balance.ValueT1 != nil && (t1Date == nil || (*balance.DatetimeEntity).After(*t1Date)) {
						t1 = balance.ValueT1
						t1Date = &balance.DatetimeEntity.Time
					}

					if balance.ValueT2 != nil && (t2Date == nil || (*balance.DatetimeEntity).After(*t2Date)) {
						t2 = balance.ValueT2
						t2Date = &balance.DatetimeEntity.Time
					}

					if balance.ValueT3 != nil && (t3Date == nil || (*balance.DatetimeEntity).After(*t3Date)) {
						t3 = balance.ValueT3
						t3Date = &balance.DatetimeEntity.Time
					}
				}

				if t1 != nil {
					val := *t1 * 1000

					metricServiceBalance.With("account", account.Number, "service", serviceID, "tariff", "1").Set(val)

					if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValueT1.Format(account.Number, serviceID), val); e != nil {
						err = multierr.Append(err, e)
					}
				}

				if t2 != nil {
					val := *t2 * 1000

					metricServiceBalance.With("account", account.Number, "service", serviceID, "tariff", "2").Set(val)

					if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValueT2.Format(account.Number, serviceID), val); e != nil {
						err = multierr.Append(err, e)
					}
				}

				if t3 != nil {
					val := *t3 * 1000

					metricServiceBalance.With("account", account.Number, "service", serviceID, "tariff", "3").Set(val)

					if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValueT3.Format(account.Number, serviceID), val); e != nil {
						err = multierr.Append(err, e)
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

package elektroset

import (
	"context"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/integratorit/elektroset"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
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

	dateStart := time.Now().Add(-b.config.BalanceDetailsInterval)
	dateEnd := time.Now()

	type value struct {
		tariff string
		value  float64
		date   time.Time
	}

	for _, account := range accounts {
		var totalBalance float64

		for _, service := range account.Services {
			totalBalance += service.Balance
			serviceID := strconv.FormatUint(service.ID, 10)

			metricServiceBalance.With("account", account.Number, "service", serviceID).Set(service.Balance)

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicServiceBalance.Format(account.Number, serviceID), service.Balance); e != nil {
				err = multierr.Append(err, e)
			}

			// balance
			if details, e := b.client.BalanceDetails(ctx, account.Number, service.Provider, dateStart, dateEnd); e == nil {
				values := make([]value, 3)
				lastBill := time.Time{}

				for _, balance := range details {
					switch balance.KDTariffPlanEntity {
					// выставленные счета
					case elektroset.TariffPlanEntityBill:
						if lastBill.IsZero() || balance.DatetimeEntity.After(lastBill) {
							lastBill = balance.DatetimeEntity.Time
						}

					// показания
					case elektroset.TariffPlanEntityValue:
						if v := values[0]; balance.ValueT1 != nil && (v.date.IsZero() || balance.DatetimeEntity.After(v.date)) {
							values[0] = value{
								tariff: "1",
								value:  *balance.ValueT1,
								date:   balance.DatetimeEntity.Time,
							}
						}

						if v := values[1]; balance.ValueT2 != nil && (v.date.IsZero() || balance.DatetimeEntity.After(v.date)) {
							values[1] = value{
								tariff: "2",
								value:  *balance.ValueT2,
								date:   balance.DatetimeEntity.Time,
							}
						}

						if v := values[2]; balance.ValueT3 != nil && (v.date.IsZero() || balance.DatetimeEntity.After(v.date)) {
							values[2] = value{
								tariff: "3",
								value:  *balance.ValueT3,
								date:   balance.DatetimeEntity.Time,
							}
						}
					}
				}

				for _, v := range values {
					if v.date.IsZero() {
						continue
					}

					v.value *= 1000

					metricMeterValue.With("account", account.Number, "service", serviceID, "tariff", v.tariff).Set(v.value)

					if e := b.MQTT().PublishAsync(ctx, b.config.TopicMeterValue.Format(account.Number, serviceID, v.tariff), v.value); e != nil {
						err = multierr.Append(err, e)
					}

					if e := b.MQTT().PublishAsync(ctx, b.config.TopicMeterDate.Format(account.Number, serviceID, v.tariff), v.date); e != nil {
						err = multierr.Append(err, e)
					}
				}

				if !lastBill.IsZero() {
					billLink, e := b.Widget().URL(map[string]string{
						"action": "bill",
						"period": lastBill.Format(layoutPeriod),
					})

					if e == nil {
						if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastBill.Format(account.Number, serviceID), billLink); e != nil {
							err = multierr.Append(err, e)
						}
					}
				}
			} else {
				err = multierr.Append(err, e)
			}
		}

		metricBalance.With("account", account.Number).Set(totalBalance)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance.Format(account.Number), totalBalance); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

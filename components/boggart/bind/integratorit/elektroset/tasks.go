package elektroset

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/performance"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	houses, err := b.Houses(ctx)
	if err != nil {
		return err
	}

	cfg := b.config()
	dateStart := time.Now().Add(-cfg.BalanceDetailsInterval)
	dateEnd := time.Now()

	type value struct {
		tariff string
		value  float64
		date   time.Time
	}

	for _, house := range houses {
		var totalBalance float64
		houseID := strconv.FormatUint(house.ID, 10)

		for _, service := range house.Services {
			totalBalance += service.Balance
			serviceID := strconv.FormatUint(service.ID, 10)

			metricServiceBalance.With("house", houseID, "service", serviceID).Set(service.Balance)

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicServiceBalance.Format(houseID, serviceID), service.Balance); e != nil {
				err = multierr.Append(err, e)
			}

			// meter checkup
			if meter, e := b.client.Meter(ctx, service.AccountID, service.Provider); e == nil {
				if e := b.MQTT().PublishAsync(ctx, cfg.TopicMeterCheckupDate.Format(houseID, serviceID), meter.CalibrationDate.Time); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}

			// balance
			if details, e := b.client.BalanceDetails(ctx, service.AccountID, service.Provider, dateStart, dateEnd); e == nil {
				values := make([]value, 3)
				lastBill := time.Time{}

				for _, balance := range details {
					switch balance.TariffPlanEntityID {
					// выставленные счета
					case elektroset.TariffPlanEntityBill:
						if lastBill.IsZero() || balance.DatetimeEntity.After(lastBill) {
							lastBill = balance.DatetimeEntity.Time
						}

						// показания
					case elektroset.TariffPlanEntityValue:
						if v := values[0]; balance.T1Value != nil && (v.date.IsZero() || balance.DatetimeEntity.After(v.date)) {
							values[0] = value{
								tariff: "1",
								value:  *balance.T1Value,
								date:   balance.DatetimeEntity.Time,
							}
						}

						if v := values[1]; balance.T2Value != nil && (v.date.IsZero() || balance.DatetimeEntity.After(v.date)) {
							values[1] = value{
								tariff: "2",
								value:  *balance.T2Value,
								date:   balance.DatetimeEntity.Time,
							}
						}

						if v := values[2]; balance.T3Value != nil && (v.date.IsZero() || balance.DatetimeEntity.After(v.date)) {
							values[2] = value{
								tariff: "3",
								value:  *balance.T3Value,
								date:   balance.DatetimeEntity.Time,
							}
						}
					}
				}

				var metersCount uint32

				for _, v := range values {
					if v.date.IsZero() {
						continue
					}

					metersCount++
					v.value *= 1000

					metricMeterValue.With("house", houseID, "service", serviceID, "tariff", v.tariff).Set(v.value)

					if e := b.MQTT().PublishAsync(ctx, cfg.TopicMeterValue.Format(houseID, serviceID, v.tariff), v.value); e != nil {
						err = multierr.Append(err, e)
					}

					if e := b.MQTT().PublishAsync(ctx, cfg.TopicMeterDate.Format(houseID, serviceID, v.tariff), v.date); e != nil {
						err = multierr.Append(err, e)
					}
				}

				b.metersCount.Set(metersCount)

				if !lastBill.IsZero() {
					provider, _ := json.Marshal(service.Provider)

					billLink, e := b.Widget().URL(map[string]string{
						"action":   "bill",
						"period":   lastBill.Format(layoutPeriod),
						"account":  service.AccountID,
						"provider": performance.UnsafeBytes2String(provider),
					})

					if e == nil {
						if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastBill.Format(houseID, serviceID), billLink); e != nil {
							err = multierr.Append(err, e)
						}
					}
				}
			} else {
				err = multierr.Append(err, e)
			}
		}

		metricBalance.With("house", houseID).Set(totalBalance)

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicBalance.Format(house.ID), totalBalance); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

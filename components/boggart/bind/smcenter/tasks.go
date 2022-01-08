package smcenter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/smcenter/client/accounting"
	"github.com/kihamo/boggart/providers/smcenter/client/meters"
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

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	cfg := b.config()

	accountingResponse, e := b.provider.Accounting.AccountingInfo(accounting.NewAccountingInfoParamsWithContext(ctx))
	if e == nil {
		for _, account := range accountingResponse.GetPayload().Data {
			metricAccountBalance.With("account", account.Ident).Set(account.TotalSum)

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicAccountBalance.Format(account.Ident), account.Sum); e != nil {
				err = fmt.Errorf("publish account balance failed: %w", e)
			}

			if len(account.Bills) > 0 {
				for _, bill := range account.Bills {
					if !bill.HasFile {
						continue
					}

					billLink, e := b.Widget().URL(map[string]string{
						"action": "bill",
						"bill":   strconv.FormatUint(bill.ID, 10),
					})

					if e == nil {
						if e := b.MQTT().PublishAsync(ctx, cfg.TopicAccountBill.Format(account.Ident), billLink); e != nil {
							err = fmt.Errorf("publish account bill failed: %w", e)
						}
					}

					break
				}
			}
		}
	} else {
		err = fmt.Errorf("get accounting info failed: %w", e)
	}

	metersResponse, e := b.provider.Meters.List(meters.NewListParamsWithContext(ctx))
	if e == nil {
		for _, meter := range metersResponse.GetPayload().Data {
			if meter.IsDisabled || len(meter.Values) == 0 {
				continue
			}

			metricMeterLastValue.With("account", meter.Ident, "factory-number", meter.FactoryNumber).Set(meter.Values[0].Value)

			var checkupDate time.Time
			if !meter.NextCheckupDate.Time().IsZero() {
				checkupDate = meter.NextCheckupDate.Time()
			} else if !meter.LastCheckupDate.Time().IsZero() && meter.RecheckInterval > 0 {
				checkupDate = meter.LastCheckupDate.Time().AddDate(int(meter.RecheckInterval), 0, 0)
			}

			if !checkupDate.IsZero() {
				if e := b.MQTT().PublishAsync(ctx, cfg.TopicMeterCheckupDate.Format(meter.Ident, meter.FactoryNumber), checkupDate); e != nil {
					err = fmt.Errorf("publish meter checkup date failed: %w", e)
				}
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicMeterValue.Format(meter.Ident, meter.FactoryNumber), meter.Values[0].Value); e != nil {
				err = fmt.Errorf("publish meter value failed: %w", e)
			}
		}
	} else {
		err = fmt.Errorf("get meters list failed: %w", e)
	}

	return err
}

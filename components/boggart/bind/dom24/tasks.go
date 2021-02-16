package dom24

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/dom24/client/accounting"
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
	accountingResponse, err := b.provider.Accounting.AccountingInfo(accounting.NewAccountingInfoParamsWithContext(ctx))
	if err != nil {
		return fmt.Errorf("get accounting info failed: %w", err)
	}

	cfg := b.config()

	for _, account := range accountingResponse.GetPayload().Data {
		metricAccountBalance.With("account", account.Ident).Set(account.TotalSum)

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicAccountBalance.Format(account.Ident), account.TotalSum); e != nil {
			err = multierr.Append(err, e)
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
						err = multierr.Append(err, e)
					}
				}

				break
			}
		}
	}

	return err
}

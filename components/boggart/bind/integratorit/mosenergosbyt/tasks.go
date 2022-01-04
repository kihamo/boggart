package mosenergosbyt

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
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
	account, err := b.Account(ctx)
	if err != nil {
		return err
	}

	cfg := b.config()

	if balance, e := b.client.CurrentBalance(ctx, account); e != nil {
		err = multierr.Append(e, err)
	} else {
		metricBalance.With("account", account.AccountID).Set(balance)

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicBalance.Format(account.AccountID), balance); e != nil {
			err = multierr.Append(e, err)
		}

		//var serviceID string
		//
		//for i, service := range balance.Services {
		//	if id, ok := services[strings.ToLower(service.Service)]; ok {
		//		serviceID = id
		//	} else {
		//		serviceID = strconv.FormatInt(int64(i), 10)
		//	}
		//
		//	metricServiceBalance.With("account", account.AccountID, "service", serviceID).Set(service.Balance)
		//
		//	if e := b.MQTT().PublishAsync(ctx, cfg.TopicServiceBalance.Format(account.AccountID, serviceID), service.Balance); e != nil {
		//		err = multierr.Append(e, err)
		//	}
		//}
	}

	// last bill

	if bills, e := b.client.Bills(ctx, account); e != nil {
		err = multierr.Append(e, err)
	} else if len(bills) > 0 {
		downloadLink, e := b.Widget().URL(map[string]string{
			"action": "download",
			"bill":   bills[0].ID,
		})

		if e == nil {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastBill.Format(account.AccountID), downloadLink); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

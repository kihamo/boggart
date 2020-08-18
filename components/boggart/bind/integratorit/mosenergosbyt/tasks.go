package mosenergosbyt

import (
	"context"
	"strconv"
	"strings"
	"time"

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
	account, err := b.Account(ctx)
	if err != nil {
		return err
	}

	if balance, e := b.client.CurrentBalance(ctx, account.Provider.IDAbonent); e != nil {
		err = multierr.Append(e, err)
	} else {
		metricBalance.With("account", account.NNAccount).Set(balance.Balance)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance.Format(account.NNAccount), balance.Balance); e != nil {
			err = multierr.Append(e, err)
		}

		for i, service := range balance.Services {
			var serviceID string

			if id, ok := services[strings.ToLower(service.Service)]; ok {
				serviceID = id
			} else {
				serviceID = strconv.FormatInt(int64(i), 10)
			}

			metricServiceBalance.With("account", account.NNAccount, "service", serviceID).Set(service.Balance)

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicServiceBalance.Format(account.NNAccount, serviceID), service.Balance); e != nil {
				err = multierr.Append(e, err)
			}
		}
	}

	// last bill
	if details, e := b.client.ChargeDetail(ctx, account.Provider.IDAbonent, time.Now().Add(-b.config.BalanceDetailsInterval), time.Now()); e != nil {
		err = multierr.Append(e, err)
	} else {
		lastBillTime := time.Time{}
		lastBillUUID := ""

		for _, detail := range details {
			if lastBillTime.IsZero() || detail.Period.After(lastBillTime) {
				lastBillTime = detail.Period.Time
				lastBillUUID = detail.Services[0].ReportUUID
			}
		}

		if !lastBillTime.IsZero() {
			billLink, e := b.Widget().URL(map[string]string{
				"action": "bill",
				"period": lastBillTime.Format(layoutPeriod),
				"uuid":   lastBillUUID,
			})

			if e == nil {
				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastBill.Format(account.NNAccount), billLink); e != nil {
					err = multierr.Append(err, e)
				}
			}
		}
	}

	return err
}

package rkcm

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/rkcm/client/general"
	"github.com/kihamo/boggart/providers/rkcm/client/mobile"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	// TODO: предположительно скидывает кэш, но надо проверить
	paramsServices := general.NewGetAdditionalServicesParamsWithContext(ctx).
		WithLogin(b.config.Login).
		WithPwd(b.config.Password)
	if _, e := b.client.General.GetAdditionalServices(paramsServices); e != nil {
		err = multierr.Append(err, e)
	}

	paramsDebt := mobile.NewGetDebtParamsWithContext(ctx).
		WithPhone(b.config.Login)
	if responseDebt, e := b.client.Mobile.GetDebt(paramsDebt); e == nil {
		for _, debt := range responseDebt.Payload.Data {
			paramsAccount := general.NewGetDebtByAccountParamsWithContext(ctx).
				WithIdent(debt.Ident)

			if responseAccount, e := b.client.General.GetDebtByAccount(paramsAccount); e == nil {
				metricBalance.With("ident", debt.Ident).Set(debt.Sum)

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance.Format(debt.Ident), responseAccount.Payload.Data.Sum); e != nil {
					err = multierr.Append(e, err)
				}
			} else {
				err = multierr.Append(e, err)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	paramsMeters := general.NewGetMeterValuesEverydayModeParamsWithContext(ctx).
		WithLogin(b.config.Login).
		WithPwd(b.config.Password)
	if responseMeters, e := b.client.General.GetMeterValuesEverydayMode(paramsMeters); e == nil {
		for _, meter := range responseMeters.Payload.Meter {
			value, e := strconv.ParseFloat(strings.ReplaceAll(meter.Value[0].Value, ",", "."), 64)
			if e != nil {
				err = multierr.Append(err, e)
				continue
			}

			metricMeterValue.With("number", strconv.FormatInt(meter.FactoryNumber, 10)).Set(value)

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicMeterValue.Format(meter.Ident, meter.FactoryNumber), value); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}

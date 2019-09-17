package rkcm

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/rkcm/client/general"
	"github.com/kihamo/boggart/providers/rkcm/client/mobile"
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
	paramsDebt := mobile.NewGetDebtParamsWithContext(ctx).
		WithPhone(b.config.Login)

	responseDebt, err := b.client.Mobile.GetDebt(paramsDebt)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
	} else {
		b.UpdateStatus(boggart.BindStatusOnline)
	}

	if err == nil {
		for _, debt := range responseDebt.Payload.Data {
			metricBalance.With("ident", debt.Ident).Set(debt.Sum)

			if e := b.MQTTPublishAsync(ctx, b.config.TopicBalance.Format(debt.Ident), debt.Sum); e != nil {
				err = multierr.Append(e, err)
			}
		}
	}

	paramsMeters := general.NewGetMeterValuesEverydayModeParamsWithContext(ctx).
		WithLogin(b.config.Login).
		WithPwd(b.config.Password)

	if responseMeters, e := b.client.General.GetMeterValuesEverydayMode(paramsMeters); e == nil {
		for _, meter := range responseMeters.Payload.Meter {
			value, e := strconv.ParseFloat(strings.Replace(meter.Value[0].Value, ",", ".", -1), 64)
			if e != nil {
				err = multierr.Append(err, e)
				continue
			}

			metricMeterValue.With("number", strconv.FormatInt(meter.FactoryNumber, 10)).Set(value)

			if e := b.MQTTPublishAsync(ctx, b.config.TopicMeterValue.Format(meter.Ident, meter.FactoryNumber), value); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	return nil, err
}

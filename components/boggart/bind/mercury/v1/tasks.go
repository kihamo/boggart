package v1

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.updaterInterval)
	taskStateUpdater.SetName("updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	t1, t2, t3, t4, err := b.provider.PowerCounters()
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)
	mTariff := metricTariff.With("serial_number", sn)

	mTariff.With("tariff", "1").Set(float64(t1))

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 1), t1); e != nil {
		err = multierr.Append(err, e)
	}

	mTariff.With("tariff", "2").Set(float64(t2))

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 2), t2); e != nil {
		err = multierr.Append(err, e)
	}

	mTariff.With("tariff", "3").Set(float64(t3))

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 3), t3); e != nil {
		err = multierr.Append(err, e)
	}

	mTariff.With("tariff", "4").Set(float64(t4))

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 4), t4); e != nil {
		err = multierr.Append(err, e)
	}

	// optimization
	if voltage, amperage, power, e := b.provider.ParamsCurrent(); e == nil {
		metricVoltage.With("serial_number", sn).Set(float64(voltage))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicVoltage.Format(snMQTT), voltage); e != nil {
			err = multierr.Append(err, e)
		}

		metricAmperage.With("serial_number", sn).Set(amperage)

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAmperage.Format(snMQTT), amperage); e != nil {
			err = multierr.Append(err, e)
		}

		metricPower.With("serial_number", sn).Set(float64(power))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicPower.Format(snMQTT), power); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if voltage, e := b.provider.BatteryVoltage(); e == nil {
		metricBatteryVoltage.With("serial_number", sn).Set(voltage)

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBatteryVoltage.Format(snMQTT), voltage); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.LastPowerOffDatetime(); e == nil {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicLastPowerOff.Format(snMQTT), date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.LastPowerOnDatetime(); e == nil {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicLastPowerOn.Format(snMQTT), date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.MakeDate(); e == nil {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicMakeDate.Format(snMQTT), date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if version, date, e := b.provider.Version(); e == nil {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicFirmwareDate.Format(snMQTT), date); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicFirmwareVersion.Format(snMQTT), version); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return nil, err
}

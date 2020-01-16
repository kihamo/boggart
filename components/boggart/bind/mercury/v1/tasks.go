package v1

import (
	"context"

	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	t1, t2, t3, t4, err := b.provider.PowerCounters()
	if err != nil {
		return err
	}

	sn := b.SerialNumber()
	mTariff := metricTariff.With("serial_number", sn)

	mTariff.With("tariff", "1").Set(float64(t1))

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicTariff1, t1); e != nil {
		err = multierr.Append(err, e)
	}

	mTariff.With("tariff", "2").Set(float64(t2))

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicTariff2, t2); e != nil {
		err = multierr.Append(err, e)
	}

	mTariff.With("tariff", "3").Set(float64(t3))

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicTariff3, t3); e != nil {
		err = multierr.Append(err, e)
	}

	mTariff.With("tariff", "4").Set(float64(t4))

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicTariff4, t4); e != nil {
		err = multierr.Append(err, e)
	}

	// optimization
	if voltage, amperage, power, e := b.provider.ParamsCurrent(); e == nil {
		metricVoltage.With("serial_number", sn).Set(float64(voltage))

		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicVoltage, voltage); e != nil {
			err = multierr.Append(err, e)
		}

		metricAmperage.With("serial_number", sn).Set(amperage)

		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAmperage, amperage); e != nil {
			err = multierr.Append(err, e)
		}

		metricPower.With("serial_number", sn).Set(float64(power))

		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicPower, power); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if voltage, e := b.provider.BatteryVoltage(); e == nil {
		metricBatteryVoltage.With("serial_number", sn).Set(voltage)

		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicBatteryVoltage, voltage); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.LastPowerOffDatetime(); e == nil {
		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicLastPowerOff, date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.LastPowerOnDatetime(); e == nil {
		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicLastPowerOn, date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.MakeDate(); e == nil {
		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicMakeDate, date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if version, date, e := b.provider.Version(); e == nil {
		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicFirmwareDate, date); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicFirmwareVersion, version); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}

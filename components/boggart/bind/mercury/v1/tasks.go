package v1

import (
	"context"

	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	provider, err := b.Provider()
	if err != nil {
		return err
	}

	tariffCount, err := provider.TariffCount()
	if err != nil {
		return err
	}

	powerValues, err := provider.PowerCounters()
	if err != nil {
		return err
	}

	sn := b.Meta().SerialNumber()
	mTariff := metricTariff.With("serial_number", sn)

	if tariffCount > 0 {
		mTariff.With("tariff", "1").Set(float64(powerValues.Tariff1()))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicTariff1, powerValues.Tariff1()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if tariffCount > 1 {
		mTariff.With("tariff", "2").Set(float64(powerValues.Tariff2()))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicTariff2, powerValues.Tariff2()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if tariffCount > 2 {
		mTariff.With("tariff", "3").Set(float64(powerValues.Tariff3()))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicTariff3, powerValues.Tariff3()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if tariffCount > 3 {
		mTariff.With("tariff", "4").Set(float64(powerValues.Tariff4()))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicTariff4, powerValues.Tariff4()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	// optimization
	if voltage, amperage, power, e := provider.UIPCurrent(); e == nil {
		metricVoltage.With("serial_number", sn).Set(float64(voltage))
		metricAmperage.With("serial_number", sn).Set(amperage)
		metricPower.With("serial_number", sn).Set(float64(power))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicVoltage, voltage); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicAmperage, amperage); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicPower, power); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if voltage, e := provider.BatteryVoltage(); e == nil {
		metricBatteryVoltage.With("serial_number", sn).Set(voltage)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBatteryVoltage, voltage); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := provider.LastPowerOffDatetime(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastPowerOff, date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := provider.LastPowerOnDatetime(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastPowerOn, date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := provider.MakeDate(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicMakeDate, date); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if version, date, e := provider.Version(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicFirmwareDate, date); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicFirmwareVersion, version); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}

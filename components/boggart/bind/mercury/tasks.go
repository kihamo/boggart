package mercury

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
	taskStateUpdater.SetName("bind-mercury:200-updater-" + b.SerialNumber())

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

	if ok := b.tariff1.Set(uint32(t1)); ok {
		mTariff.With("tariff", "1").Set(float64(t1))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 1), 0, true, t1); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.tariff2.Set(uint32(t2)); ok {
		mTariff.With("tariff", "2").Set(float64(t2))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 2), 0, true, t2); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.tariff3.Set(uint32(t3)); ok {
		mTariff.With("tariff", "3").Set(float64(t3))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 3), 0, true, t3); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.tariff4.Set(uint32(t4)); ok {
		mTariff.With("tariff", "4").Set(float64(t4))

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicTariff.Format(snMQTT, 4), 0, true, t4); e != nil {
			err = multierr.Append(err, e)
		}
	}

	// optimization
	if voltage, amperage, power, e := b.provider.ParamsCurrent(); e == nil {
		if ok := b.voltage.Set(uint32(voltage)); ok {
			metricVoltage.With("serial_number", sn).Set(float64(voltage))

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicVoltage.Format(snMQTT), 0, true, voltage); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if ok := b.amperage.Set(float32(amperage)); ok {
			metricAmperage.With("serial_number", sn).Set(amperage)

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAmperage.Format(snMQTT), 0, true, amperage); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if ok := b.power.Set(uint32(power)); ok {
			metricPower.With("serial_number", sn).Set(float64(power))

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicPower.Format(snMQTT), 0, true, power); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if voltage, e := b.provider.BatteryVoltage(); e == nil {
		if ok := b.batteryVoltage.Set(float32(voltage)); ok {
			metricBatteryVoltage.With("serial_number", sn).Set(voltage)

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBatteryVoltage.Format(snMQTT), 0, true, voltage); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.LastPowerOffDatetime(); e == nil {
		if ok := b.lastPowerOffDate.Set(uint32(date.Unix())); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicLastPowerOff.Format(snMQTT), 0, true, date); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.LastPowerOnDatetime(); e == nil {
		if ok := b.lastPowerOnDate.Set(uint32(date.Unix())); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicLastPowerOn.Format(snMQTT), 0, true, date); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if date, e := b.provider.MakeDate(); e == nil {
		if ok := b.makeDate.Set(uint32(date.Unix())); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicMakeDate.Format(snMQTT), 0, true, date); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if version, date, e := b.provider.Version(); e == nil {
		if ok := b.firmwareDate.Set(uint32(date.Unix())); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicFirmwareDate.Format(snMQTT), 0, true, date); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if ok := b.firmwareVersion.Set(version); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicFirmwareVersion.Format(snMQTT), 0, true, version); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	return nil, err
}

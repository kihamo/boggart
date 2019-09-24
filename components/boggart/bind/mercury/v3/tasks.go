package v3

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/mercury/v3"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (_ interface{}, err error) {
	if err = b.provider.ChannelTest(); err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	sn := b.SerialNumber()
	if sn == "" {
		var (
			makeDate time.Time
			version  string
		)

		sn, makeDate, version, _, err = b.provider.ForceReadParameters()
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
			return nil, err
		}

		b.SetSerialNumber(sn)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicMakeDate.Format(sn), makeDate); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicFirmwareVersion.Format(sn), version); e != nil {
			err = multierr.Append(err, e)
		}
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	if val, _, _, _, e := b.provider.ReadArray(v3.ArrayReset, nil, v3.Tariff1); e == nil {
		metricTariff.With("serial_number", sn).With("tariff", "1").Set(float64(val))

		if e := b.MQTTPublishAsync(ctx, b.config.TopicTariff1.Format(sn), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if p1, p2, p3, e := b.provider.Voltage(); e == nil {
		metricVoltage.With("serial_number", sn).With("phase", "1").Set(p1)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicVoltage1.Format(sn), p1); e != nil {
			err = multierr.Append(err, e)
		}

		metricVoltage.With("serial_number", sn).With("phase", "2").Set(p2)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicVoltage2.Format(sn), p2); e != nil {
			err = multierr.Append(err, e)
		}

		metricVoltage.With("serial_number", sn).With("phase", "3").Set(p3)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicVoltage3.Format(sn), p3); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if p1, p2, p3, e := b.provider.Amperage(); e == nil {
		metricAmperage.With("serial_number", sn).With("phase", "1").Set(p1)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicAmperage1.Format(sn), p1); e != nil {
			err = multierr.Append(err, e)
		}

		metricAmperage.With("serial_number", sn).With("phase", "2").Set(p2)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicAmperage2.Format(sn), p2); e != nil {
			err = multierr.Append(err, e)
		}

		metricAmperage.With("serial_number", sn).With("phase", "3").Set(p3)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicAmperage3.Format(sn), p3); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if _, p1, p2, p3, e := b.provider.Power(v3.PowerNumberP); e == nil {
		metricPower.With("serial_number", sn).With("phase", "1").Set(p1)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicPower1.Format(sn), p1); e != nil {
			err = multierr.Append(err, e)
		}

		metricPower.With("serial_number", sn).With("phase", "2").Set(p2)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicPower2.Format(sn), p2); e != nil {
			err = multierr.Append(err, e)
		}

		metricPower.With("serial_number", sn).With("phase", "3").Set(p3)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicPower3.Format(sn), p3); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return nil, err
}
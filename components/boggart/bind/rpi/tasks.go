package rpi

import (
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/rpi"
	"go.uber.org/multierr"
)

var voltsIDs = []rpi.VoltsID{
	rpi.VoltsIDCore,
	rpi.VoltsIDSDramC,
	rpi.VoltsIDSDramI,
	rpi.VoltsIDSDramP,
}

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
	if value, e := b.providerSysFS.Model(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicModel, value); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	sn := b.serialNumber

	if values, e := b.providerSysFS.CPUFrequentie(); e == nil {
		for num, value := range values {
			metricCPUFrequentie.With("serial_number", sn).With("cpu", "cpu"+strconv.FormatUint(num, 10)).Set(float64(value))

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicCPUFrequentie.Format(num), value); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if value, e := b.providerSysFS.Temperature(); e == nil {
		metricTemperature.With("serial_number", sn).Set(value)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicTemperature, value); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if value, e := b.providerSysFS.Throttled(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicCurrentlyUnderVoltage, value.IsCurrentlyUnderVoltage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicCurrentlyThrottled, value.IsCurrentlyThrottled()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicCurrentlyARMFrequencyCapped, value.IsCurrentlyARMFrequencyCapped()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicCurrentlySoftTemperatureReached, value.IsCurrentlySoftTemperatureReached()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicSinceRebootUnderVoltage, value.IsSinceRebootUnderVoltage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicSinceRebootThrottled, value.IsSinceRebootThrottled()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicSinceRebootARMFrequencyCapped, value.IsSinceRebootARMFrequencyCapped()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicSinceRebootSoftTemperatureReached, value.IsSinceRebootSoftTemperatureReached()); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	for _, id := range voltsIDs {
		if value, e := b.providerVCGenCMD.Voltage(id); e == nil {
			metricVoltage.With("serial_number", sn).With("id", id.String()).Set(value)

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicVoltage.Format(id), value); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	return err
}

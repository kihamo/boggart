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
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()
	cfg := b.config()

	if value, e := b.providerSysFS.Model(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicModel.Format(sn), value); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if values, e := b.providerSysFS.CPUFrequentie(); e == nil {
		for num, value := range values {
			metricCPUFrequentie.With("serial_number", sn).With("cpu", "cpu"+strconv.FormatUint(num, 10)).Set(float64(value))

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicCPUFrequentie.Format(sn, num), value); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if value, e := b.providerSysFS.Temperature(); e == nil {
		metricTemperature.With("serial_number", sn).Set(value)

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicTemperature.Format(sn), value); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if value, e := b.providerSysFS.Throttled(); e == nil {
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentlyUnderVoltage.Format(sn), value.IsCurrentlyUnderVoltage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentlyThrottled.Format(sn), value.IsCurrentlyThrottled()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentlyARMFrequencyCapped.Format(sn), value.IsCurrentlyARMFrequencyCapped()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentlySoftTemperatureReached.Format(sn), value.IsCurrentlySoftTemperatureReached()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicSinceRebootUnderVoltage.Format(sn), value.IsSinceRebootUnderVoltage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicSinceRebootThrottled.Format(sn), value.IsSinceRebootThrottled()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicSinceRebootARMFrequencyCapped.Format(sn), value.IsSinceRebootARMFrequencyCapped()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicSinceRebootSoftTemperatureReached.Format(sn), value.IsSinceRebootSoftTemperatureReached()); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	for _, id := range voltsIDs {
		if value, e := b.providerVCGenCMD.Voltage(id); e == nil {
			metricVoltage.With("serial_number", sn).With("id", id.String()).Set(value)

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicVoltage.Format(sn, id), value); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	return err
}

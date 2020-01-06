package rpi

import (
	"context"

	"github.com/kihamo/boggart/providers/rpi"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

var voltsIDs = []rpi.VoltsID{
	rpi.VoltsIDCore,
	rpi.VoltsIDSDramC,
	rpi.VoltsIDSDramI,
	rpi.VoltsIDSDramP,
}

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (err error) {
	if value, e := b.providerSysFS.Model(); e == nil {
		if e := b.MQTTPublishAsync(ctx, b.config.TopicModel, value); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if values, e := b.providerSysFS.CPUFrequentie(); e == nil {
		for num, value := range values {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicCPUFrequentie.Format(num), value); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if value, e := b.providerSysFS.Temperature(); e == nil {
		if e := b.MQTTPublishAsync(ctx, b.config.TopicTemperature, value); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if value, e := b.providerSysFS.Throttled(); e == nil {
		if e := b.MQTTPublishAsync(ctx, b.config.TopicCurrentlyUnderVoltage, value.IsCurrentlyUnderVoltage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicCurrentlyThrottled, value.IsCurrentlyThrottled()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicCurrentlyARMFrequencyCapped, value.IsCurrentlyARMFrequencyCapped()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicCurrentlySoftTemperatureReached, value.IsCurrentlySoftTemperatureReached()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicSinceRebootUnderVoltage, value.IsSinceRebootUnderVoltage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicSinceRebootThrottled, value.IsSinceRebootThrottled()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicSinceRebootARMFrequencyCapped, value.IsSinceRebootARMFrequencyCapped()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicSinceRebootSoftTemperatureReached, value.IsSinceRebootSoftTemperatureReached()); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	for _, id := range voltsIDs {
		if value, e := b.providerVCGenCMD.Voltage(id); e == nil {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicVoltage.Format(id), value); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	return err
}

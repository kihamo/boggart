package sun

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	b.taskStateUpdater = task.NewFunctionTask(b.taskUpdater)
	b.taskStateUpdater.SetRepeats(-1)
	b.taskStateUpdater.SetName("bind-astro-sun-updater-" + b.SerialNumber())

	return []workers.Task{
		b.taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (_ interface{}, err error) {
	serialNumberMQTT := mqtt.NameReplace(b.SerialNumber())
	times := b.Times()

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNadir.Format(serialNumberMQTT), times.Nadir); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightBeforeStart.Format(serialNumberMQTT), times.NightBefore.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightBeforeEnd.Format(serialNumberMQTT), times.NightBefore.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightBeforeDuration.Format(serialNumberMQTT), times.NightBefore.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAstronomicalDawnStart.Format(serialNumberMQTT), times.AstronomicalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAstronomicalDawnEnd.Format(serialNumberMQTT), times.AstronomicalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAstronomicalDawnDuration.Format(serialNumberMQTT), times.AstronomicalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNauticalDawnStart.Format(serialNumberMQTT), times.NauticalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNauticalDawnEnd.Format(serialNumberMQTT), times.NauticalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNauticalDawnDuration.Format(serialNumberMQTT), times.NauticalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCivilDawnStart.Format(serialNumberMQTT), times.CivilDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCivilDawnEnd.Format(serialNumberMQTT), times.CivilDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCivilDawnDuration.Format(serialNumberMQTT), times.CivilDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRiseStart.Format(serialNumberMQTT), times.Sunrise.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRiseEnd.Format(serialNumberMQTT), times.Sunrise.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRiseDuration.Format(serialNumberMQTT), times.Sunrise.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSolarNoon.Format(serialNumberMQTT), times.SolarNoon); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSetStart.Format(serialNumberMQTT), times.Sunset.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSetEnd.Format(serialNumberMQTT), times.Sunset.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSetDuration.Format(serialNumberMQTT), times.Sunset.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCivilDuskStart.Format(serialNumberMQTT), times.CivilDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCivilDuskEnd.Format(serialNumberMQTT), times.CivilDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCivilDuskDuration.Format(serialNumberMQTT), times.CivilDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNauticalDuskStart.Format(serialNumberMQTT), times.NauticalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNauticalDuskEnd.Format(serialNumberMQTT), times.NauticalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNauticalDuskDuration.Format(serialNumberMQTT), times.NauticalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAstronomicalDuskStart.Format(serialNumberMQTT), times.AstronomicalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAstronomicalDuskEnd.Format(serialNumberMQTT), times.AstronomicalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicAstronomicalDuskDuration.Format(serialNumberMQTT), times.AstronomicalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightAfterStart.Format(serialNumberMQTT), times.NightAfter.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightAfterEnd.Format(serialNumberMQTT), times.NightAfter.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightAfterDuration.Format(serialNumberMQTT), times.NightAfter.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	// change start date (only one of 24 hours)
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	b.taskStateUpdater.SetRepeatInterval(todayStart.Add(dayDuration).Sub(now))

	return nil, err
}

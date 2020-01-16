package sun

import (
	"context"
	"time"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	b.taskStateUpdater = task.NewFunctionTask(b.taskUpdater)
	b.taskStateUpdater.SetRepeats(-1)
	b.taskStateUpdater.SetName("updater")

	return []workers.Task{
		b.taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (_ interface{}, err error) {
	times := b.Times()

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNadir, times.Nadir); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNightBeforeStart, times.NightBefore.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNightBeforeEnd, times.NightBefore.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNightBeforeDuration, times.NightBefore.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAstronomicalDawnStart, times.AstronomicalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAstronomicalDawnEnd, times.AstronomicalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAstronomicalDawnDuration, times.AstronomicalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNauticalDawnStart, times.NauticalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNauticalDawnEnd, times.NauticalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNauticalDawnDuration, times.NauticalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicCivilDawnStart, times.CivilDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicCivilDawnEnd, times.CivilDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicCivilDawnDuration, times.CivilDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicRiseStart, times.Sunrise.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicRiseEnd, times.Sunrise.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicRiseDuration, times.Sunrise.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicSolarNoon, times.SolarNoon); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicSetStart, times.Sunset.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicSetEnd, times.Sunset.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicSetDuration, times.Sunset.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicCivilDuskStart, times.CivilDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicCivilDuskEnd, times.CivilDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicCivilDuskDuration, times.CivilDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNauticalDuskStart, times.NauticalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNauticalDuskEnd, times.NauticalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNauticalDuskDuration, times.NauticalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAstronomicalDuskStart, times.AstronomicalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAstronomicalDuskEnd, times.AstronomicalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicAstronomicalDuskDuration, times.AstronomicalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNightAfterStart, times.NightAfter.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNightAfterEnd, times.NightAfter.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTContainer().PublishAsync(ctx, b.config.TopicNightAfterDuration, times.NightAfter.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	// change start date (only one of 24 hours)
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	b.taskStateUpdater.SetRepeatInterval(todayStart.Add(dayDuration).Sub(now))

	return nil, err
}

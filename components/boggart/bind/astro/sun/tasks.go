package sun

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithSchedule(tasks.ScheduleWithDailyTime(tasks.ScheduleNow(), 0, 0, 0, nil)).
			WithHandlerFunc(b.taskUpdaterHandler),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	times := b.Times()

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNadir, times.Nadir); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNightBeforeStart, times.NightBefore.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNightBeforeEnd, times.NightBefore.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNightBeforeDuration, times.NightBefore.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicAstronomicalDawnStart, times.AstronomicalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicAstronomicalDawnEnd, times.AstronomicalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicAstronomicalDawnDuration, times.AstronomicalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNauticalDawnStart, times.NauticalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNauticalDawnEnd, times.NauticalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNauticalDawnDuration, times.NauticalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCivilDawnStart, times.CivilDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCivilDawnEnd, times.CivilDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCivilDawnDuration, times.CivilDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicRiseStart, times.Sunrise.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicRiseEnd, times.Sunrise.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicRiseDuration, times.Sunrise.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicSolarNoon, times.SolarNoon); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicSetStart, times.Sunset.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicSetEnd, times.Sunset.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicSetDuration, times.Sunset.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCivilDuskStart, times.CivilDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCivilDuskEnd, times.CivilDusk.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicCivilDuskDuration, times.CivilDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNauticalDuskStart, times.NauticalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNauticalDuskEnd, times.NauticalDusk.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNauticalDuskDuration, times.NauticalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicAstronomicalDuskStart, times.AstronomicalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicAstronomicalDuskEnd, times.AstronomicalDusk.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicAstronomicalDuskDuration, times.AstronomicalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNightAfterStart, times.NightAfter.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNightAfterEnd, times.NightAfter.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicNightAfterDuration, times.NightAfter.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

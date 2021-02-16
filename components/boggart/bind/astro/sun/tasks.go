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
	cfg := b.config()
	id := b.Meta().ID()

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNadir.Format(id), times.Nadir); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNightBeforeStart.Format(id), times.NightBefore.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNightBeforeEnd.Format(id), times.NightBefore.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNightBeforeDuration.Format(id), times.NightBefore.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicAstronomicalDawnStart.Format(id), times.AstronomicalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicAstronomicalDawnEnd.Format(id), times.AstronomicalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicAstronomicalDawnDuration.Format(id), times.AstronomicalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNauticalDawnStart.Format(id), times.NauticalDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNauticalDawnEnd.Format(id), times.NauticalDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNauticalDawnDuration.Format(id), times.NauticalDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCivilDawnStart.Format(id), times.CivilDawn.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCivilDawnEnd.Format(id), times.CivilDawn.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCivilDawnDuration.Format(id), times.CivilDawn.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicRiseStart.Format(id), times.Sunrise.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicRiseEnd.Format(id), times.Sunrise.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicRiseDuration.Format(id), times.Sunrise.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicSolarNoon.Format(id), times.SolarNoon); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicSetStart.Format(id), times.Sunset.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicSetEnd.Format(id), times.Sunset.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicSetDuration.Format(id), times.Sunset.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCivilDuskStart.Format(id), times.CivilDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCivilDuskEnd.Format(id), times.CivilDusk.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCivilDuskDuration.Format(id), times.CivilDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNauticalDuskStart.Format(id), times.NauticalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNauticalDuskEnd.Format(id), times.NauticalDusk.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNauticalDuskDuration.Format(id), times.NauticalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicAstronomicalDuskStart.Format(id), times.AstronomicalDusk.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicAstronomicalDuskEnd.Format(id), times.AstronomicalDusk.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicAstronomicalDuskDuration.Format(id), times.AstronomicalDusk.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNightAfterStart.Format(id), times.NightAfter.Start); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNightAfterEnd.Format(id), times.NightAfter.End); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicNightAfterDuration.Format(id), times.NightAfter.Duration); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

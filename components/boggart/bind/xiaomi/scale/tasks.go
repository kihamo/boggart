package scale

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
	"go.uber.org/multierr"
)

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

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	measures, err := b.Measures(ctx)
	if err != nil {
		return err
	}

	mac := b.Meta().MACAsString()
	cfg := b.config()
	profile := b.CurrentProfile()
	measureStartDatetime := b.measureStartDatetime.Load()

	for _, measure := range measures {
		b.Logger().Debug("Measure",
			"start_datetime", measureStartDatetime.String(),
			"datetime", measure.Datetime().String(),
			"unit", measure.Unit(),
			"weight", measure.Weight(),
			"impedance", measure.Impedance(),
			"profile", profile.Name,
		)

		dt := measure.Datetime()

		// если метрика снята после до установки профиля, то она может относится к другому профилю
		// и испортит показатели текущего профиля
		if !dt.After(measureStartDatetime) {
			continue
		}

		// в v2 impedance равен 0 в промежуточных результах взвешивания,
		// поэтому такое значение можно игнорироть
		// TODO: сделать настраиваемо это поведение
		if measure.Impedance() == 0 || measure.Impedance() > 3000 {
			continue
		}

		metricWeight.With("mac", mac, "profile", profile.Name).Set(measure.Weight())
		metricImpedance.With("mac", mac, "profile", profile.Name).Set(float64(measure.Impedance()))

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicDatetime.Format(mac, profile.Name), measure.Datetime()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicWeight.Format(mac, profile.Name), measure.Weight()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicImpedance.Format(mac, profile.Name), measure.Impedance()); e != nil {
			err = multierr.Append(err, e)
		}

		// skip guest profile
		if profile.Height == 0 || profile.GetAge() == 0 {
			if err == nil && dt.After(measureStartDatetime) {
				b.measureStartDatetime.Set(dt)
				measureStartDatetime = dt
			}

			continue
		}

		sex := scale.SexMale
		if profile.Sex {
			sex = scale.SexFemale
		}

		metrics, e := measure.Metrics(sex, profile.Height, profile.GetAge())
		if e != nil {
			err = multierr.Append(err, e)
			continue
		}

		// ever
		metricBMR.With("mac", mac, "profile", profile.Name).Set(float64(metrics.BMR()))
		metricBMR.With("mac", mac, "profile", profile.Name).Set(metrics.BMI())
		metricFatPercentage.With("mac", mac, "profile", profile.Name).Set(metrics.FatPercentage())
		metricWaterPercentage.With("mac", mac, "profile", profile.Name).Set(metrics.WaterPercentage())
		metricIdealWeight.With("mac", mac, "profile", profile.Name).Set(metrics.IdealWeight())

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicBMR.Format(mac, profile.Name), metrics.BMR()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicBMI.Format(mac, profile.Name), metrics.BMI()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicFatPercentage.Format(mac, profile.Name), metrics.FatPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicWaterPercentage.Format(mac, profile.Name), metrics.WaterPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicIdealWeight.Format(mac, profile.Name), metrics.IdealWeight()); e != nil {
			err = multierr.Append(err, e)
		}

		// only sets impedance
		if measure.Impedance() > 0 {
			metricLBMCoefficient.With("mac", mac, "profile", profile.Name).Set(metrics.LBMCoefficient())
			metricBoneMass.With("mac", mac, "profile", profile.Name).Set(metrics.BoneMass())
			metricMuscleMass.With("mac", mac, "profile", profile.Name).Set(metrics.MuscleMass())
			metricVisceralFat.With("mac", mac, "profile", profile.Name).Set(metrics.VisceralFat())
			metricFatMassToIdeal.With("mac", mac, "profile", profile.Name).Set(metrics.FatMassToIdeal())
			metricProteinPercentage.With("mac", mac, "profile", profile.Name).Set(metrics.ProteinPercentage())
			metricBodyType.With("mac", mac, "profile", profile.Name).Set(float64(metrics.BodyType()))
			metricMetabolicAge.With("mac", mac, "profile", profile.Name).Set(metrics.MetabolicAge())

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicLBMCoefficient.Format(mac, profile.Name), metrics.LBMCoefficient()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicBoneMass.Format(mac, profile.Name), metrics.BoneMass()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicMuscleMass.Format(mac, profile.Name), metrics.MuscleMass()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicVisceralFat.Format(mac, profile.Name), metrics.VisceralFat()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicFatMassToIdeal.Format(mac, profile.Name), metrics.FatMassToIdeal()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicProteinPercentage.Format(mac, profile.Name), metrics.ProteinPercentage()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicBodyType.Format(mac, profile.Name), metrics.BodyType()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicMetabolicAge.Format(mac, profile.Name), metrics.MetabolicAge()); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if err == nil && dt.After(measureStartDatetime) {
			b.measureStartDatetime.Set(dt)
			measureStartDatetime = dt
		}
	}

	return err
}

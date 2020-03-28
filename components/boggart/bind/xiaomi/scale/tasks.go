package scale

import (
	"context"

	"github.com/kihamo/boggart/providers/xiaomi/scale"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	measures, err := b.Measures(ctx)
	if err != nil {
		return err
	}

	mac := b.Meta().MACAsString()
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

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicDatetime.Format(profile.Name), measure.Datetime()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicWeight.Format(profile.Name), measure.Weight()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicImpedance.Format(profile.Name), measure.Impedance()); e != nil {
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

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicBMR.Format(profile.Name), metrics.BMR()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicBMI.Format(profile.Name), metrics.BMI()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicFatPercentage.Format(profile.Name), metrics.FatPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicWaterPercentage.Format(profile.Name), metrics.WaterPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicIdealWeight.Format(profile.Name), metrics.IdealWeight()); e != nil {
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

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicLBMCoefficient.Format(profile.Name), metrics.LBMCoefficient()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicBoneMass.Format(profile.Name), metrics.BoneMass()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicMuscleMass.Format(profile.Name), metrics.MuscleMass()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicVisceralFat.Format(profile.Name), metrics.VisceralFat()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicFatMassToIdeal.Format(profile.Name), metrics.FatMassToIdeal()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicProteinPercentage.Format(profile.Name), metrics.ProteinPercentage()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicBodyType.Format(profile.Name), metrics.BodyType()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicMetabolicAge.Format(profile.Name), metrics.MetabolicAge()); e != nil {
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

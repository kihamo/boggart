package scale

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/xiaomi/scale"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	profile := b.CurrentProfile()
	if profile == nil {
		return errors.New("profile isn't set")
	}

	measures, err := b.provider.Measures(ctx)
	if err != nil {
		return err
	}

	sn := b.SerialNumber()
	setDatetime := b.setProfileDatetime.Load()

	for _, measure := range measures {
		// если метрика снята после до установки профиля, то она может относится к другому профилю
		// и испортит показатели текущего профиля
		if measure.Datetime().Before(setDatetime) {
			continue
		}

		// в v2 impedance равен 0 в промежуточных результах взвешивания,
		// поэтому такое значение можно игнорироть
		// TODO: сделать настраиваемо это поведение
		if measure.Impedance() == 0 {
			continue
		}

		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicDatetime.Format(profile.Name), measure.Datetime()); e != nil {
			err = multierr.Append(err, e)
		}

		metricWeight.With("serial_number", sn, "profile", profile.Name).Set(measure.Weight())
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicWeight.Format(profile.Name), measure.Weight()); e != nil {
			err = multierr.Append(err, e)
		}

		metricImpedance.With("serial_number", sn, "profile", profile.Name).Set(float64(measure.Impedance()))
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicImpedance.Format(profile.Name), measure.Impedance()); e != nil {
			err = multierr.Append(err, e)
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
		metricBMR.With("serial_number", sn, "profile", profile.Name).Set(float64(metrics.BMR()))
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBMR.Format(profile.Name), metrics.BMR()); e != nil {
			err = multierr.Append(err, e)
		}

		metricBMR.With("serial_number", sn, "profile", profile.Name).Set(metrics.BMI())
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBMI.Format(profile.Name), metrics.BMI()); e != nil {
			err = multierr.Append(err, e)
		}

		metricFatPercentage.With("serial_number", sn, "profile", profile.Name).Set(metrics.FatPercentage())
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicFatPercentage.Format(profile.Name), metrics.FatPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		metricWaterPercentage.With("serial_number", sn, "profile", profile.Name).Set(metrics.WaterPercentage())
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicWaterPercentage.Format(profile.Name), metrics.WaterPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		metricIdealWeight.With("serial_number", sn, "profile", profile.Name).Set(metrics.IdealWeight())
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicIdealWeight.Format(profile.Name), metrics.IdealWeight()); e != nil {
			err = multierr.Append(err, e)
		}

		// only sets impedance
		if measure.Impedance() > 0 {
			metricLBMCoefficient.With("serial_number", sn, "profile", profile.Name).Set(metrics.LBMCoefficient())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicLBMCoefficient.Format(profile.Name), metrics.LBMCoefficient()); e != nil {
				err = multierr.Append(err, e)
			}

			metricBoneMass.With("serial_number", sn, "profile", profile.Name).Set(metrics.BoneMass())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBoneMass.Format(profile.Name), metrics.BoneMass()); e != nil {
				err = multierr.Append(err, e)
			}

			metricMuscleMass.With("serial_number", sn, "profile", profile.Name).Set(metrics.MuscleMass())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicMuscleMass.Format(profile.Name), metrics.MuscleMass()); e != nil {
				err = multierr.Append(err, e)
			}

			metricVisceralFat.With("serial_number", sn, "profile", profile.Name).Set(metrics.VisceralFat())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicVisceralFat.Format(profile.Name), metrics.VisceralFat()); e != nil {
				err = multierr.Append(err, e)
			}

			metricFatMassToIdeal.With("serial_number", sn, "profile", profile.Name).Set(metrics.FatMassToIdeal())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicFatMassToIdeal.Format(profile.Name), metrics.FatMassToIdeal()); e != nil {
				err = multierr.Append(err, e)
			}

			metricProteinPercentage.With("serial_number", sn, "profile", profile.Name).Set(metrics.ProteinPercentage())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicProteinPercentage.Format(profile.Name), metrics.ProteinPercentage()); e != nil {
				err = multierr.Append(err, e)
			}

			metricBodyType.With("serial_number", sn, "profile", profile.Name).Set(float64(metrics.BodyType()))
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBodyType.Format(profile.Name), metrics.BodyType()); e != nil {
				err = multierr.Append(err, e)
			}

			metricMetabolicAge.With("serial_number", sn, "profile", profile.Name).Set(metrics.MetabolicAge())
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicMetabolicAge.Format(profile.Name), metrics.MetabolicAge()); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

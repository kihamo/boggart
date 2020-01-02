package scale

import (
	"context"

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
	measures, err := b.provider.Measures(ctx)
	if err != nil {
		return err
	}

	for _, measure := range measures {
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicDatetime, measure.Datetime()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicWeight, measure.Weight()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicImpedance, measure.Impedance()); e != nil {
			err = multierr.Append(err, e)
		}

		if b.sex.IsNil() || b.height.IsNil() || b.age.IsNil() {
			continue
		}

		sex := scale.SexMale
		if b.sex.IsTrue() {
			sex = scale.SexFemale
		}

		metrics, e := measure.Metrics(sex, uint64(b.height.Load()), uint64(b.age.Load()))
		if e != nil {
			err = multierr.Append(err, e)
			continue
		}

		// ever
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBMR, metrics.BMR()); e != nil {
			err = multierr.Append(err, e)
		}
		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBMI, metrics.BMI()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicFatPercentage, metrics.FatPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicWaterPercentage, metrics.WaterPercentage()); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicIdealWeight, metrics.IdealWeight()); e != nil {
			err = multierr.Append(err, e)
		}

		// only sets impedance
		if measure.Impedance() > 0 {
			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicLBMCoefficient, metrics.LBMCoefficient()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBoneMass, metrics.BoneMass()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicMuscleMass, metrics.MuscleMass()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicVisceralFat, metrics.VisceralFat()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicFatMassToIdeal, metrics.FatMassToIdeal()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicProteinPercentage, metrics.ProteinPercentage()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicBodyType, metrics.BodyType()); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTTPublishAsyncWithoutCache(ctx, b.config.TopicMetabolicAge, metrics.MetabolicAge()); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

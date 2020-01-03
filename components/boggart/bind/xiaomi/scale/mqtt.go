package scale

import (
	"context"
	"encoding/json"
	"errors"

	"go.uber.org/multierr"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicProfile,
		b.config.TopicDatetime,
		b.config.TopicWeight,
		b.config.TopicImpedance,
		b.config.TopicBMR,
		b.config.TopicBMI,
		b.config.TopicFatPercentage,
		b.config.TopicWaterPercentage,
		b.config.TopicIdealWeight,
		b.config.TopicLBMCoefficient,
		b.config.TopicBoneMass,
		b.config.TopicMuscleMass,
		b.config.TopicVisceralFat,
		b.config.TopicFatMassToIdeal,
		b.config.TopicProteinPercentage,
		b.config.TopicBodyType,
		b.config.TopicMetabolicAge,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicProfileActivate, 0, b.callbackMQTTProfileActive),
	}
}

func (b *Bind) callbackMQTTProfileActive(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	if len(route) < 2 {
		return errors.New("bad topic name")
	}

	b.SetProfile(route[len(route)-2])
	return b.notifyCurrentProfile(ctx)
}

func (b *Bind) notifyCurrentProfile(ctx context.Context) error {
	profile := b.CurrentProfile()
	if profile == nil {
		return errors.New("profile isn't set")
	}

	response, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicProfile, response); e != nil {
		err = multierr.Append(err, e)
	}

	return nil
}

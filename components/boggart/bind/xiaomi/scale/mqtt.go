package scale

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
	"go.uber.org/multierr"
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
	response, err := json.Marshal(b.CurrentProfile())
	if err != nil {
		return err
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicProfile, response); e != nil {
		err = multierr.Append(err, e)
	}

	return nil
}

package scale

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
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
		mqtt.NewSubscriber(b.config.TopicProfile, 0, b.callbackMQTTProfile),
	}
}

func (b *Bind) callbackMQTTProfile(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	var request struct {
		Sex    bool   `json:"sex"`
		Height uint32 `json:"height"`
		Age    uint32 `json:"age"`
	}

	if err := message.UnmarshalJSON(&request); err != nil {
		return err
	}

	if request.Height <= 0 || request.Height > 220 {
		return errors.New("height is either too low or too high (limits: <0cm and >220cm)")
	}

	if request.Age <= 0 || request.Age > 99 {
		return errors.New("age is either too low or too high (limits: <= 0 years and > 99 years)")
	}

	b.sex.Set(request.Sex)
	b.height.Set(request.Height)
	b.age.Set(request.Age)

	return nil
}

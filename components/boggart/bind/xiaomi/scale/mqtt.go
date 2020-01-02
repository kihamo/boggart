package scale

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.uber.org/multierr"

	"github.com/kihamo/boggart/components/mqtt"
)

type profile struct {
	Sex      *bool      `json:"sex,omitempty"`
	Height   *uint32    `json:"height,omitempty"`
	Birthday *time.Time `json:"birthday,omitempty"`
	Age      *uint32    `json:"age,omitempty"`
}

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
		mqtt.NewSubscriber(b.config.TopicProfileSet, 0, b.callbackMQTTProfile),
	}
}

func (b *Bind) callbackMQTTProfile(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	var request profile

	if err := message.UnmarshalJSON(&request); err != nil {
		return err
	}

	if request.Sex == nil {
		return errors.New("sex isn't set")
	}

	if request.Height == nil {
		return errors.New("height isn't set")
	}

	if request.Age == nil && request.Birthday == nil {
		return errors.New("age or birthday isn't set")
	}

	sex := *request.Sex
	height := *request.Height
	age := uint32(0)

	if height <= 0 || height > 220 {
		return errors.New("height is either too low or too high (limits: <0cm and >220cm)")
	}

	if request.Age != nil {
		age = *request.Age
	} else if request.Birthday != nil && !request.Birthday.IsZero() {
		now := time.Now()

		age = uint32(now.Year() - request.Birthday.Year())

		diff := time.Date(now.Year(), request.Birthday.Month(), request.Birthday.Day(), request.Birthday.Hour(), request.Birthday.Minute(), request.Birthday.Second(), request.Birthday.Nanosecond(), request.Birthday.Location())
		if diff.After(now) {
			age -= 1
		}
	}

	if age <= 0 || age > 99 {
		return errors.New("age is either too low or too high (limits: <= 0 years and > 99 years)")
	}

	b.sex.Set(sex)
	b.height.Set(height)
	b.age.Set(age)

	b.updateProfile(ctx)

	return nil
}

func (b *Bind) updateProfile(ctx context.Context) error {
	response := profile{}

	if !b.sex.IsNil() {
		response.Sex = &[]bool{b.sex.Load()}[0]
	}

	if !b.height.IsNil() {
		response.Height = &[]uint32{b.height.Load()}[0]
	}

	if !b.age.IsNil() {
		response.Age = &[]uint32{b.age.Load()}[0]
	}

	value, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicProfile, value); e != nil {
		err = multierr.Append(err, e)
	}

	return nil
}

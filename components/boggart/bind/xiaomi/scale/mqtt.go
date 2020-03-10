package scale

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
	"go.uber.org/multierr"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicProfileActivate, 0, b.callbackMQTTProfileActive),
		mqtt.NewSubscriber(b.config.TopicProfileSettings, 0, b.callbackMQTTProfileSettings),
	}
}

func (b *Bind) callbackMQTTProfileActive(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	if len(route) < 2 {
		return errors.New("bad topic name")
	}

	b.SetProfile(route[len(route)-2])
	return b.notifyCurrentProfile(ctx)
}

func (b *Bind) callbackMQTTProfileSettings(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	if len(route) < 2 {
		return errors.New("bad topic name")
	}

	var settings Profile

	if err := json.Unmarshal(message.Payload(), &settings); err != nil {
		return err
	}

	profileGuest.Sex = settings.Sex
	profileGuest.Height = settings.Height
	profileGuest.Birthday = settings.Birthday
	profileGuest.Age = settings.Age

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

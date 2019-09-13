package service

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicOnline,
		b.config.TopicLatency,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicCheck, 0, func(ctx context.Context, _ mqtt.Component, _ mqtt.Message) error {
			_, err := b.taskUpdater(ctx)
			return err
		}),
	}
}

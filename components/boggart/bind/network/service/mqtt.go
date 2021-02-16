package service

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config().TopicCheck.Format(b.address), 0, func(ctx context.Context, _ mqtt.Component, _ mqtt.Message) error {
			return b.ReadinessProbe(ctx)
		}),
	}
}

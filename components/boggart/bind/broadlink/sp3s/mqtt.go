package sp3s

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicState,
		b.config.TopicPower,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicSet, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			if message.IsTrue() {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
	}
}

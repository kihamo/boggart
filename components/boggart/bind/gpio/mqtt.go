package gpio

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicPinState,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if b.Mode() != ModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			b.config.TopicPinSet,
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if message.IsTrue() {
					return b.High(ctx)
				}

				return b.Low(ctx)
			}),
	}
}

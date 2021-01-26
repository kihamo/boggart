package gpio

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if b.Mode() != ModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			b.config.TopicPinSet,
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if message.IsTrue() && !b.config.Inverted {
					return b.High(ctx)
				}

				return b.Low(ctx)
			}),
	}
}

func (b *Bind) publishState(ctx context.Context, value bool) error {
	if b.config.Inverted {
		value = !value
	}

	return b.MQTT().PublishAsync(ctx, b.config.TopicPinState, value)
}

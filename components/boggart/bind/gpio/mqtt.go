package gpio

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if b.Mode() != ModeOut {
		return nil
	}

	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			cfg.TopicPinSet.Format(b.pin.Number()),
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if message.IsTrue() && !cfg.Inverted {
					return b.High(ctx)
				}

				return b.Low(ctx)
			}),
	}
}

func (b *Bind) publishState(ctx context.Context, value bool) error {
	cfg := b.config()

	if cfg.Inverted {
		value = !value
	}

	return b.MQTT().PublishAsync(ctx, cfg.TopicPinState.Format(b.pin.Number()), value)
}

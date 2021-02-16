package sp3s

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicSet.Format(cfg.MAC.String()), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			if message.IsTrue() {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
	}
}

package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			b.config.TopicShutdown,
			0,
			func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				if message.IsTrue() {
					return b.application.Shutdown()
				}

				return nil
			}),
	}
}

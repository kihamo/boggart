package fcm

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicSend, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Send(ctx, message.String())
		})),
	}
}

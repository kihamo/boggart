package fcm

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/fcm/+"

	MQTTSubscribeTopicSendMessage = MQTTPrefix + "/send"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicSendMessage.Format(b.projectID), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Send(ctx, message.String())
		}),
	}
}

package sp3s

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/socket/+/"

	MQTTPublishTopicState = MQTTPrefix + "state"
	MQTTPublishTopicPower = MQTTPrefix + "power"
	MQTTSubscribeTopicSet = MQTTPrefix + "set"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicState.Format(sn)),
		mqtt.Topic(MQTTPublishTopicPower.Format(sn)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicSet.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			if message.IsTrue() {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
	}
}

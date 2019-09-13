package gpio

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/gpio/+"

	MQTTPublishTopicPinState = MQTTPrefix
	MQTTSubscribeTopicPinSet = MQTTPrefix + "/set"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicPinState.Format(b.pin.Number()),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if b.Mode() != ModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			MQTTSubscribeTopicPinSet.Format(b.pin.Number()),
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if message.IsTrue() {
					return b.High(ctx)
				}

				return b.Low(ctx)
			}),
	}
}

package gpio

import (
	"bytes"
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicPinState mqtt.Topic = boggart.ComponentName + "/gpio/+"
	MQTTSubscribeTopicPinSet mqtt.Topic = boggart.ComponentName + "/gpio/+/set"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicPinState.Format(b.pin.Number())),
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
				if bytes.Equal(message.Payload(), []byte(`1`)) {
					return b.High(ctx)
				}

				return b.Low(ctx)
			}),
	}
}

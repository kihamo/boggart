package broadlink

import (
	"bytes"
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	SP3SMQTTPublishTopicState mqtt.Topic = boggart.ComponentName + "/socket/+/state"
	SP3SMQTTPublishTopicPower mqtt.Topic = boggart.ComponentName + "/socket/+/power"
	SP3SMQTTSubscribeTopicSet mqtt.Topic = boggart.ComponentName + "/socket/+/set"
)

func (b *BindSP3S) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(SP3SMQTTPublishTopicState.Format(sn)),
		mqtt.Topic(SP3SMQTTPublishTopicPower.Format(sn)),
	}
}

func (b *BindSP3S) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(SP3SMQTTSubscribeTopicSet.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			if bytes.Equal(message.Payload(), []byte(`1`)) {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
	}
}

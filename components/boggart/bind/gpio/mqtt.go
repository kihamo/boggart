package gpio

import (
	"bytes"
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicPinState mqtt.Topic = boggart.ComponentName + "/gpio/+"
	MQTTTopicPinSet   mqtt.Topic = boggart.ComponentName + "/gpio/+/set"
)

func (d *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicPinState.Format(d.pin.Number())),
	}
}

func (d *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if d.Mode() != ModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			MQTTTopicPinSet.Format(d.pin.Number()),
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
				if bytes.Equal(message.Payload(), []byte(`1`)) {
					return d.High(ctx)
				}

				return d.Low(ctx)
			}),
	}
}

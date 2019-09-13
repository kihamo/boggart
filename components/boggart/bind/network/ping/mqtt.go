package ping

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/ping/+/"

	MQTTPublishTopicOnline  = MQTTPrefix + "online"
	MQTTPublishTopicLatency = MQTTPrefix + "latency"
	MQTTSubscribeTopicCheck = MQTTPrefix + "check"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicLatency.Format(b.hostname),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicCheck.Format(b.hostname), 0, func(ctx context.Context, _ mqtt.Component, _ mqtt.Message) error {
			_, err := b.taskUpdater(ctx)
			return err
		}),
	}
}

package network

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	PingMQTTPrefix mqtt.Topic = boggart.ComponentName + "/ping/+/"

	PingMQTTPublishTopicOnline  = PingMQTTPrefix + "online"
	PingMQTTPublishTopicLatency = PingMQTTPrefix + "latency"
	PingMQTTSubscribeTopicCheck = PingMQTTPrefix + "check"
)

func (b *BindPing) MQTTPublishes() []mqtt.Topic {
	h := mqtt.NameReplace(b.hostname)

	return []mqtt.Topic{
		mqtt.Topic(PingMQTTPublishTopicLatency.Format(h)),
	}
}

func (b *BindPing) MQTTSubscribers() []mqtt.Subscriber {
	h := mqtt.NameReplace(b.hostname)

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(PingMQTTSubscribeTopicCheck.Format(h), 0, func(ctx context.Context, _ mqtt.Component, _ mqtt.Message) error {
			_, err := b.taskUpdater(ctx)
			return err
		}),
	}
}

package service

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/service/+/"

	MQTTPublishTopicOnline  = MQTTPrefix + "online"
	MQTTPublishTopicLatency = MQTTPrefix + "latency"
	MQTTSubscribeTopicCheck = MQTTPrefix + "check"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	h := mqtt.NameReplace(b.address)

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicLatency.Format(h)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	h := mqtt.NameReplace(b.address)

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicCheck.Format(h), 0, func(ctx context.Context, _ mqtt.Component, _ mqtt.Message) error {
			_, err := b.taskUpdater(ctx)
			return err
		}),
	}
}

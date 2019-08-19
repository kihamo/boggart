package xmeye

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicEvent mqtt.Topic = boggart.ComponentName + "/cctv/+/event/+/+"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := make([]mqtt.Topic, 0, 1)

	if b.config.AlarmStreamingEnabled {
		topics = append(topics, MQTTPublishTopicEvent)
	}

	return topics
}

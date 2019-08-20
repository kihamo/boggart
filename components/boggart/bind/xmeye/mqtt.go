package xmeye

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicEvent            mqtt.Topic = boggart.ComponentName + "/cctv/+/event/+/+"
	MQTTPublishTopicStateHDDCapacity mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/capacity"
	MQTTPublishTopicStateHDDFree     mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/free"
	MQTTPublishTopicStateHDDUsage    mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/usage"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
		MQTTPublishTopicStateHDDCapacity,
		MQTTPublishTopicStateHDDFree,
		MQTTPublishTopicStateHDDUsage,
	}

	if b.config.AlarmStreamingEnabled {
		topics = append(topics, MQTTPublishTopicEvent)
	}

	return topics
}

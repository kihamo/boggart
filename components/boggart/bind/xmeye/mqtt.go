package xmeye

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
		b.config.TopicStateModel,
		b.config.TopicStateFirmwareVersion,
		b.config.TopicStateFirmwareReleasedDate,
		b.config.TopicStateHDDCapacity,
		b.config.TopicStateHDDFree,
		b.config.TopicStateHDDUsage,
	}

	if b.config.AlarmStreamingEnabled {
		topics = append(topics, b.config.TopicEvent)
	}

	return topics
}

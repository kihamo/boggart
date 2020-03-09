package tvt

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
		b.config.TopicStateModel,
		b.config.TopicStateFirmwareVersion,
		b.config.TopicStateHDDCapacity,
		b.config.TopicStateHDDFree,
		b.config.TopicStateHDDUsage,
	}

	return topics
}

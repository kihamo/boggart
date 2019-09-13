package boggart

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicName,
		b.config.TopicVersion,
		b.config.TopicBuild,
	}
}

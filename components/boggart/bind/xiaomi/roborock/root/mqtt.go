package root

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicRuntimeConfig mqtt.Topic = boggart.ComponentName + "/+/runtime/+"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicRuntimeConfig.Format(b.SerialNumber()),
	}
}

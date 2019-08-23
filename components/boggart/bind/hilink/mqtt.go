package hilink

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicSMS mqtt.Topic = boggart.ComponentName + "/hilink/+/sms"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicSMS,
	}
}

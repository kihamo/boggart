package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicBalance mqtt.Topic = boggart.ComponentName + "/service/softvideo/+/balance"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicBalance.Format(b.SerialNumber()),
	}
}

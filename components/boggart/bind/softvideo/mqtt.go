package softvideo

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicBalance mqtt.Topic = boggart.ComponentName + "/service/softvideo/+/balance"
)

func (b *Bind) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicBalance.Format(b.SerialNumber())),
	}
}

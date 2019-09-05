package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicBalance mqtt.Topic = boggart.ComponentName + "/service/mosenergosbyt/+/balance"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicBalance,
	}
}

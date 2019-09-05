package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/service/elektroset/+/"

	MQTTPublishTopicBalance        = MQTTPrefix + "balance"
	MQTTPublishTopicServiceBalance = MQTTPrefix + "+/balance"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicBalance,
		MQTTPublishTopicServiceBalance,
	}
}

package elektroset

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicBalance,
		b.config.TopicServiceBalance,
		b.config.TopicMeterValueT1,
		b.config.TopicMeterValueT2,
		b.config.TopicMeterValueT3,
	}
}

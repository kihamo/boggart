package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicValue mqtt.Topic = boggart.ComponentName + "/meter/ds18b20/+"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicValue.Format(mqtt.NameReplace(b.SerialNumber()))),
	}
}

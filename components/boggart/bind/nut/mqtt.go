package nut

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicVariable mqtt.Topic = boggart.ComponentName + "/ups/+/+"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicVariable.Format(b.config.UPS)),
	}
}

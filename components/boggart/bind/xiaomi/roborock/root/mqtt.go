package root

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicRuntimeConfig mqtt.Topic = boggart.ComponentName + "/+/runtime/+"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicRuntimeConfig.Format(sn)),
	}
}

func (b *Bind) SetMQTTClient(client mqtt.Component) {
	b.BindMQTT.SetMQTTClient(client)

	if client != nil {
		for fileName, callback := range b.watchFiles {
			go callback(fileName)
		}
	}
}

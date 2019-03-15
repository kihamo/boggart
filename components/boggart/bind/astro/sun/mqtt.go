package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/astro/sun/+/"

	MQTTPublishTopicSunrise  = MQTTPrefix + "sunrise"
	MQTTPublishTopicSunset   = MQTTPrefix + "sunset"
	MQTTPublishTopicDayLight = MQTTPrefix + "daylight"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicSunrise.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSunset.Format(sn)),
		mqtt.Topic(MQTTPublishTopicDayLight.Format(sn)),
	}
}

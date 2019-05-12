package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/xiaomi/roborock/+/"

	MQTTPublishTopicBattery   = MQTTPrefix + "battery"
	MQTTPublishTopicCleanArea = MQTTPrefix + "clean/area"
	MQTTPublishTopicCleanTime = MQTTPrefix + "clean/time"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicBattery,
		MQTTPublishTopicCleanArea,
		MQTTPublishTopicCleanTime,
	}
}

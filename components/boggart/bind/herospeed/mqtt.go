package herospeed

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicStateModel           mqtt.Topic = boggart.ComponentName + "/cctv/+/state/model"
	MQTTPublishTopicStateFirmwareVersion mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/version"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicStateModel,
		MQTTPublishTopicStateFirmwareVersion,
	}
}

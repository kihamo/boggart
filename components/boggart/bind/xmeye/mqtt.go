package xmeye

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicEvent                     mqtt.Topic = boggart.ComponentName + "/cctv/+/event/+/+"
	MQTTPublishTopicStateModel                mqtt.Topic = boggart.ComponentName + "/cctv/+/state/model"
	MQTTPublishTopicStateFirmwareVersion      mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/version"
	MQTTPublishTopicStateFirmwareReleasedDate mqtt.Topic = boggart.ComponentName + "/cctv/+/state/firmware/release-date"
	MQTTPublishTopicStateHDDCapacity          mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/capacity"
	MQTTPublishTopicStateHDDFree              mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/free"
	MQTTPublishTopicStateHDDUsage             mqtt.Topic = boggart.ComponentName + "/cctv/+/state/hdd/+/usage"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
		MQTTPublishTopicStateModel,
		MQTTPublishTopicStateFirmwareVersion,
		MQTTPublishTopicStateFirmwareReleasedDate,
		MQTTPublishTopicStateHDDCapacity,
		MQTTPublishTopicStateHDDFree,
		MQTTPublishTopicStateHDDUsage,
	}

	if b.config.AlarmStreamingEnabled {
		topics = append(topics, MQTTPublishTopicEvent)
	}

	return topics
}

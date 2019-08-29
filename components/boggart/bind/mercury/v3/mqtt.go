package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/"

	MQTTPublishTopicTariff          = MQTTPrefix + "tariff/+"
	MQTTPublishTopicVoltage         = MQTTPrefix + "voltage/+"
	MQTTPublishTopicAmperage        = MQTTPrefix + "amperage/+"
	MQTTPublishTopicPower           = MQTTPrefix + "power/+"
	MQTTPublishTopicMakeDate        = MQTTPrefix + "make-date"
	MQTTPublishTopicFirmwareVersion = MQTTPrefix + "firmware/version"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicTariff,
		MQTTPublishTopicVoltage,
		MQTTPublishTopicAmperage,
		MQTTPublishTopicPower,
		MQTTPublishTopicMakeDate,
		MQTTPublishTopicFirmwareVersion,
	}
}

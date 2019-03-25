package homie

import (
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix     mqtt.Topic = "+/+/"
	MQTTPrefixImpl            = MQTTPrefix + "$implementation/"

	MQTTPublishTopicBroadcast mqtt.Topic = "+/$broadcast/+"
	MQTTPublishTopicReset                = MQTTPrefixImpl + "reset"
	MQTTPublishTopicRestart              = MQTTPrefixImpl + "restart"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	base := b.config.BaseTopic
	sn := b.SerialNumber()

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicBroadcast.Format(base)),
		mqtt.Topic(MQTTPublishTopicReset.Format(base, sn)),
		mqtt.Topic(MQTTPublishTopicRestart.Format(base, sn)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	base := b.config.BaseTopic
	sn := b.SerialNumber()

	return []mqtt.Subscriber{
		// device
		mqtt.NewSubscriber(deviceMQTTSubscribeTopicAttribute.Format(base, sn), 0, b.deviceAttributesSubscriber),
		mqtt.NewSubscriber(deviceMQTTSubscribeTopicAttributeFirmware.Format(base, sn), 0, b.deviceFirmwareSubscriber),
		mqtt.NewSubscriber(deviceMQTTSubscribeTopicAttributeImplementation.Format(base, sn), 0, b.deviceImplementationSubscriber),
		mqtt.NewSubscriber(deviceMQTTSubscribeTopicAttributeStats.Format(base, sn), 0, b.deviceStatsSubscriber),

		// ota
		mqtt.NewSubscriber(otaMQTTPublishTopicStatus.Format(base, sn), 0, b.otaStatusSubscriber),

		// settings
		mqtt.NewSubscriber(settingsMQTTPublishTopicGet.Format(base, sn), 0, b.settingsSubscriber),
	}
}

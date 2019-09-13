package esp

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
		MQTTPublishTopicBroadcast.Format(base),
		MQTTPublishTopicReset.Format(base, sn),
		MQTTPublishTopicRestart.Format(base, sn),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	base := b.config.BaseTopic
	sn := b.SerialNumber()

	return []mqtt.Subscriber{
		// settings
		mqtt.NewSubscriber(settingsTopicGet.Format(base, sn), 0, b.settingsSubscriber),

		// device
		mqtt.NewSubscriber(deviceTopicAttribute.Format(base, sn), 0, b.deviceAttributesSubscriber),
		mqtt.NewSubscriber(deviceTopicAttributeFirmware.Format(base, sn), 0, b.deviceFirmwareSubscriber),
		mqtt.NewSubscriber(deviceTopicAttributeImplementation.Format(base, sn), 0, b.deviceImplementationSubscriber),
		mqtt.NewSubscriber(deviceTopicAttributeStats.Format(base, sn), 0, b.deviceStatsSubscriber),

		// nodes
		mqtt.NewSubscriber(nodesTopicNodesAttribute.Format(base, sn), 0, b.nodesAttributesSubscriber),
		mqtt.NewSubscriber(nodesTopicNodes.Format(base, sn), 0, b.nodesSubscriber),
		mqtt.NewSubscriber(nodesTopicProperty.Format(base, sn), 0, b.nodesPropertySubscriber),

		// ota
		mqtt.NewSubscriber(otaTopicStatus.Format(base, sn), 0, b.otaStatusSubscriber),
		mqtt.NewSubscriber(otaTopicEnabled.Format(base, sn), 0, b.otaEnabledSubscriber),
	}
}

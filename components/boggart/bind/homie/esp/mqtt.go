package esp

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicBroadcast,
		b.config.TopicReset,
		b.config.TopicRestart,
		b.config.TopicSettingsSet,
		b.config.TopicOTAFirmware,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		// settings
		mqtt.NewSubscriber(b.config.TopicSettings, 0, b.settingsSubscriber),

		// device
		mqtt.NewSubscriber(b.config.TopicDeviceAttribute, 0, b.deviceAttributesSubscriber),
		mqtt.NewSubscriber(b.config.TopicDeviceAttributeFirmware, 0, b.deviceFirmwareSubscriber),
		mqtt.NewSubscriber(b.config.TopicDeviceAttributeImplementation, 0, b.deviceImplementationSubscriber),
		mqtt.NewSubscriber(b.config.TopicDeviceAttributeStats, 0, b.deviceStatsSubscriber),

		// nodes
		mqtt.NewSubscriber(b.config.TopicNodeAttribute, 0, b.nodesAttributesSubscriber),
		mqtt.NewSubscriber(b.config.TopicNodeList, 0, b.nodesSubscriber),
		mqtt.NewSubscriber(b.config.TopicNodeProperty, 0, b.nodesPropertySubscriber),

		// ota
		mqtt.NewSubscriber(b.config.TopicOTAStatus, 0, b.otaStatusSubscriber),
		mqtt.NewSubscriber(b.config.TopicOTAEnabled, 0, b.otaEnabledSubscriber),
	}
}

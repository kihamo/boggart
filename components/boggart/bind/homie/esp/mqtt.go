package esp

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		// settings
		mqtt.NewSubscriber(cfg.TopicSettings, 0, b.settingsSubscriber),

		// device
		mqtt.NewSubscriber(cfg.TopicDeviceAttribute.Format(id), 0, b.deviceAttributesSubscriber),
		mqtt.NewSubscriber(cfg.TopicDeviceAttributeFirmware.Format(id), 0, b.deviceFirmwareSubscriber),
		mqtt.NewSubscriber(cfg.TopicDeviceAttributeImplementation.Format(id), 0, b.deviceImplementationSubscriber),
		mqtt.NewSubscriber(cfg.TopicDeviceAttributeStats.Format(id), 0, b.deviceStatsSubscriber),

		// nodes
		mqtt.NewSubscriber(cfg.TopicNodeAttribute.Format(id), 0, b.nodesAttributesSubscriber),
		mqtt.NewSubscriber(cfg.TopicNodeList.Format(id), 0, b.nodesSubscriber),
		mqtt.NewSubscriber(cfg.TopicNodeProperty.Format(id), 0, b.nodesPropertySubscriber),

		// ota
		mqtt.NewSubscriber(cfg.TopicOTAStatus.Format(id), 0, b.otaStatusSubscriber),
		mqtt.NewSubscriber(cfg.TopicOTAEnabled.Format(id), 0, b.otaEnabledSubscriber),
	}
}

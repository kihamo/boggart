package v1

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/"

	MQTTPublishTopicTariff          = MQTTPrefix + "tariff/+"
	MQTTPublishTopicVoltage         = MQTTPrefix + "voltage"
	MQTTPublishTopicAmperage        = MQTTPrefix + "amperage"
	MQTTPublishTopicPower           = MQTTPrefix + "power"
	MQTTPublishTopicBatteryVoltage  = MQTTPrefix + "battery_voltage"
	MQTTPublishTopicLastPowerOff    = MQTTPrefix + "last-power-off"
	MQTTPublishTopicLastPowerOn     = MQTTPrefix + "last-power-on"
	MQTTPublishTopicMakeDate        = MQTTPrefix + "make-date"
	MQTTPublishTopicFirmwareDate    = MQTTPrefix + "firmware/date"
	MQTTPublishTopicFirmwareVersion = MQTTPrefix + "firmware/version"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := b.SerialNumber()

	return []mqtt.Topic{
		MQTTPublishTopicTariff.Format(sn, 1),
		MQTTPublishTopicTariff.Format(sn, 2),
		MQTTPublishTopicTariff.Format(sn, 3),
		MQTTPublishTopicTariff.Format(sn, 4),
		MQTTPublishTopicVoltage.Format(sn),
		MQTTPublishTopicAmperage.Format(sn),
		MQTTPublishTopicPower.Format(sn),
		MQTTPublishTopicBatteryVoltage.Format(sn),
		MQTTPublishTopicLastPowerOff.Format(sn),
		MQTTPublishTopicLastPowerOn.Format(sn),
		MQTTPublishTopicMakeDate.Format(sn),
		MQTTPublishTopicFirmwareDate.Format(sn),
		MQTTPublishTopicFirmwareVersion.Format(sn),
	}
}

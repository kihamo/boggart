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
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicTariff.Format(sn, 1)),
		mqtt.Topic(MQTTPublishTopicTariff.Format(sn, 2)),
		mqtt.Topic(MQTTPublishTopicTariff.Format(sn, 3)),
		mqtt.Topic(MQTTPublishTopicTariff.Format(sn, 4)),
		mqtt.Topic(MQTTPublishTopicVoltage.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAmperage.Format(sn)),
		mqtt.Topic(MQTTPublishTopicPower.Format(sn)),
		mqtt.Topic(MQTTPublishTopicBatteryVoltage.Format(sn)),
		mqtt.Topic(MQTTPublishTopicLastPowerOff.Format(sn)),
		mqtt.Topic(MQTTPublishTopicLastPowerOn.Format(sn)),
		mqtt.Topic(MQTTPublishTopicMakeDate.Format(sn)),
		mqtt.Topic(MQTTPublishTopicFirmwareDate.Format(sn)),
		mqtt.Topic(MQTTPublishTopicFirmwareVersion.Format(sn)),
	}
}

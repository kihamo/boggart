package mercury

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicTariff          mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/tariff/+"
	MQTTPublishTopicVoltage         mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/voltage"
	MQTTPublishTopicAmperage        mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/amperage"
	MQTTPublishTopicPower           mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/power"
	MQTTPublishTopicBatteryVoltage  mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/battery_voltage"
	MQTTPublishTopicLastPowerOff    mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/last-power-off"
	MQTTPublishTopicLastPowerOn     mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/last-power-on"
	MQTTPublishTopicMakeDate        mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/make-date"
	MQTTPublishTopicFirmwareDate    mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/firmware/date"
	MQTTPublishTopicFirmwareVersion mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/firmware/version"
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

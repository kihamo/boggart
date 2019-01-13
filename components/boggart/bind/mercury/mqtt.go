package mercury

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicTariff          mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/tariff/+"
	MQTTTopicVoltage         mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/voltage"
	MQTTTopicAmperage        mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/amperage"
	MQTTTopicPower           mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/power"
	MQTTTopicBatteryVoltage  mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/battery_voltage"
	MQTTTopicLastPowerOff    mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/last-power-off"
	MQTTTopicLastPowerOn     mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/last-power-on"
	MQTTTopicMakeDate        mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/make-date"
	MQTTTopicFirmwareDate    mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/firmware/date"
	MQTTTopicFirmwareVersion mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/firmware/version"
)

func (d *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(d.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicTariff.Format(sn, 1)),
		mqtt.Topic(MQTTTopicTariff.Format(sn, 2)),
		mqtt.Topic(MQTTTopicTariff.Format(sn, 3)),
		mqtt.Topic(MQTTTopicTariff.Format(sn, 4)),
		mqtt.Topic(MQTTTopicVoltage.Format(sn)),
		mqtt.Topic(MQTTTopicAmperage.Format(sn)),
		mqtt.Topic(MQTTTopicPower.Format(sn)),
		mqtt.Topic(MQTTTopicBatteryVoltage.Format(sn)),
		mqtt.Topic(MQTTTopicLastPowerOff.Format(sn)),
		mqtt.Topic(MQTTTopicLastPowerOn.Format(sn)),
		mqtt.Topic(MQTTTopicMakeDate.Format(sn)),
		mqtt.Topic(MQTTTopicFirmwareDate.Format(sn)),
		mqtt.Topic(MQTTTopicFirmwareVersion.Format(sn)),
	}
}

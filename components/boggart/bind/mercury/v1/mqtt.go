package v1

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicTariff1,
		b.config.TopicTariff2,
		b.config.TopicTariff3,
		b.config.TopicTariff4,
		b.config.TopicVoltage,
		b.config.TopicAmperage,
		b.config.TopicPower,
		b.config.TopicBatteryVoltage,
		b.config.TopicLastPowerOff,
		b.config.TopicLastPowerOn,
		b.config.TopicMakeDate,
		b.config.TopicFirmwareDate,
		b.config.TopicFirmwareVersion,
	}
}

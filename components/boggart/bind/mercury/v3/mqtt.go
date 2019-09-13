package v3

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicTariff1,
		b.config.TopicVoltage1,
		b.config.TopicVoltage2,
		b.config.TopicVoltage3,
		b.config.TopicAmperage1,
		b.config.TopicAmperage2,
		b.config.TopicAmperage3,
		b.config.TopicPower1,
		b.config.TopicPower2,
		b.config.TopicPower3,
		b.config.TopicMakeDate,
		b.config.TopicFirmwareVersion,
	}
}

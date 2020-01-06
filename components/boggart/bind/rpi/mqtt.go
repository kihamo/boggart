package rpi

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicModel,
		b.config.TopicCPUFrequentie,
		b.config.TopicTemperature,
		b.config.TopicVoltage,
		b.config.TopicCurrentlyUnderVoltage,
		b.config.TopicCurrentlyThrottled,
		b.config.TopicCurrentlyARMFrequencyCapped,
		b.config.TopicCurrentlySoftTemperatureReached,
		b.config.TopicSinceRebootUnderVoltage,
		b.config.TopicSinceRebootThrottled,
		b.config.TopicSinceRebootARMFrequencyCapped,
		b.config.TopicSinceRebootSoftTemperatureReached,
	}
}

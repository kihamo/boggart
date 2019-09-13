package pulsar

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicTemperatureIn,
		b.config.TopicTemperatureOut,
		b.config.TopicTemperatureDelta,
		b.config.TopicEnergy,
		b.config.TopicConsumption,
		b.config.TopicCapacity,
		b.config.TopicPower,
		b.config.TopicInputPulses1,
		b.config.TopicInputPulses2,
		b.config.TopicInputPulses3,
		b.config.TopicInputPulses4,
		b.config.TopicInputVolume1,
		b.config.TopicInputVolume2,
		b.config.TopicInputVolume3,
		b.config.TopicInputVolume4,
	}
}

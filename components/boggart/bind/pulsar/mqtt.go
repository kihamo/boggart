package pulsar

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/"

	MQTTPublishTopicTemperatureIn    = MQTTPrefix + "temperature_in"
	MQTTPublishTopicTemperatureOut   = MQTTPrefix + "temperature_out"
	MQTTPublishTopicTemperatureDelta = MQTTPrefix + "temperature_delta"
	MQTTPublishTopicEnergy           = MQTTPrefix + "energy"
	MQTTPublishTopicConsumption      = MQTTPrefix + "consumption"
	MQTTPublishTopicCapacity         = MQTTPrefix + "capacity"
	MQTTPublishTopicPower            = MQTTPrefix + "power"
	MQTTPublishTopicInputPulses      = MQTTPrefix + "input/+/pulses"
	MQTTPublishTopicInputVolume      = MQTTPrefix + "input/+/volume"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := b.SerialNumber()

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicTemperatureIn.Format(sn)),
		mqtt.Topic(MQTTPublishTopicTemperatureOut.Format(sn)),
		mqtt.Topic(MQTTPublishTopicTemperatureDelta.Format(sn)),
		mqtt.Topic(MQTTPublishTopicEnergy.Format(sn)),
		mqtt.Topic(MQTTPublishTopicConsumption.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCapacity.Format(sn)),
		mqtt.Topic(MQTTPublishTopicPower.Format(sn)),
		mqtt.Topic(MQTTPublishTopicInputPulses.Format(sn, 1)),
		mqtt.Topic(MQTTPublishTopicInputVolume.Format(sn, 1)),
		mqtt.Topic(MQTTPublishTopicInputPulses.Format(sn, 2)),
		mqtt.Topic(MQTTPublishTopicInputVolume.Format(sn, 2)),
		mqtt.Topic(MQTTPublishTopicInputPulses.Format(sn, 3)),
		mqtt.Topic(MQTTPublishTopicInputVolume.Format(sn, 3)),
		mqtt.Topic(MQTTPublishTopicInputPulses.Format(sn, 4)),
		mqtt.Topic(MQTTPublishTopicInputVolume.Format(sn, 4)),
	}
}

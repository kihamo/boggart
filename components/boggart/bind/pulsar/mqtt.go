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
		MQTTPublishTopicTemperatureIn.Format(sn),
		MQTTPublishTopicTemperatureOut.Format(sn),
		MQTTPublishTopicTemperatureDelta.Format(sn),
		MQTTPublishTopicEnergy.Format(sn),
		MQTTPublishTopicConsumption.Format(sn),
		MQTTPublishTopicCapacity.Format(sn),
		MQTTPublishTopicPower.Format(sn),
		MQTTPublishTopicInputPulses.Format(sn, 1),
		MQTTPublishTopicInputVolume.Format(sn, 1),
		MQTTPublishTopicInputPulses.Format(sn, 2),
		MQTTPublishTopicInputVolume.Format(sn, 2),
		MQTTPublishTopicInputPulses.Format(sn, 3),
		MQTTPublishTopicInputVolume.Format(sn, 3),
		MQTTPublishTopicInputPulses.Format(sn, 4),
		MQTTPublishTopicInputVolume.Format(sn, 4),
	}
}

package pulsar

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicTemperatureIn    mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_in"
	MQTTPublishTopicTemperatureOut   mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_out"
	MQTTPublishTopicTemperatureDelta mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_delta"
	MQTTPublishTopicEnergy           mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/energy"
	MQTTPublishTopicConsumption      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/consumption"
	MQTTPublishTopicCapacity         mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/capacity"
	MQTTPublishTopicPower            mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/power"
	MQTTPublishTopicInputPulses      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/input/+/pulses"
	MQTTPublishTopicInputVolume      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/input/+/volume"
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

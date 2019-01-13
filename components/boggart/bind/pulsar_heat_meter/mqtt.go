package pulsar_heat_meter

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicTemperatureIn    mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_in"
	MQTTTopicTemperatureOut   mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_out"
	MQTTTopicTemperatureDelta mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_delta"
	MQTTTopicEnergy           mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/energy"
	MQTTTopicConsumption      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/consumption"
	MQTTTopicCapacity         mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/capacity"
	MQTTTopicPower            mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/power"
	MQTTTopicInputPulses      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/input/+/pulses"
	MQTTTopicInputVolume      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/input/+/volume"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := b.SerialNumber()

	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicTemperatureIn.Format(sn)),
		mqtt.Topic(MQTTTopicTemperatureOut.Format(sn)),
		mqtt.Topic(MQTTTopicTemperatureDelta.Format(sn)),
		mqtt.Topic(MQTTTopicEnergy.Format(sn)),
		mqtt.Topic(MQTTTopicConsumption.Format(sn)),
		mqtt.Topic(MQTTTopicCapacity.Format(sn)),
		mqtt.Topic(MQTTTopicPower.Format(sn)),
		mqtt.Topic(MQTTTopicInputPulses.Format(sn, 1)),
		mqtt.Topic(MQTTTopicInputVolume.Format(sn, 1)),
		mqtt.Topic(MQTTTopicInputPulses.Format(sn, 2)),
		mqtt.Topic(MQTTTopicInputVolume.Format(sn, 2)),
		mqtt.Topic(MQTTTopicInputPulses.Format(sn, 3)),
		mqtt.Topic(MQTTTopicInputVolume.Format(sn, 3)),
		mqtt.Topic(MQTTTopicInputPulses.Format(sn, 4)),
		mqtt.Topic(MQTTTopicInputVolume.Format(sn, 4)),
	}
}

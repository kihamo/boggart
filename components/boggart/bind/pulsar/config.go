package pulsar

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`

	ConnectionDSN         string `mapstructure:"connection_dsn" yaml:"connection_dsn" valid:"required"`
	Address               string
	Location              string
	Input1Offset          float32       `mapstructure:"input1_offset" valid:"float"`
	Input2Offset          float32       `mapstructure:"input2_offset" valid:"float"`
	Input3Offset          float32       `mapstructure:"input3_offset" valid:"float"`
	Input4Offset          float32       `mapstructure:"input4_offset" valid:"float"`
	UpdaterInterval       time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicTemperatureIn    mqtt.Topic    `mapstructure:"topic_temperature_in" yaml:"topic_temperature_in"`
	TopicTemperatureOut   mqtt.Topic    `mapstructure:"topic_temperature_out" yaml:"topic_temperature_out"`
	TopicTemperatureDelta mqtt.Topic    `mapstructure:"topic_temperature_delta" yaml:"topic_temperature_delta"`
	TopicEnergy           mqtt.Topic    `mapstructure:"topic_energy" yaml:"topic_energy"`
	TopicConsumption      mqtt.Topic    `mapstructure:"topic_consumption" yaml:"topic_consumption"`
	TopicCapacity         mqtt.Topic    `mapstructure:"topic_capacity" yaml:"topic_capacity"`
	TopicPower            mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicInputPulses1     mqtt.Topic    `mapstructure:"topic_pulses_1" yaml:"topic_pulses_1"`
	TopicInputPulses2     mqtt.Topic    `mapstructure:"topic_pulses_2" yaml:"topic_pulses_2"`
	TopicInputPulses3     mqtt.Topic    `mapstructure:"topic_pulses_3" yaml:"topic_pulses_3"`
	TopicInputPulses4     mqtt.Topic    `mapstructure:"topic_pulses_4" yaml:"topic_pulses_4"`
	TopicInputVolume1     mqtt.Topic    `mapstructure:"topic_volume_1" yaml:"topic_volume_1"`
	TopicInputVolume2     mqtt.Topic    `mapstructure:"topic_volume_2" yaml:"topic_volume_2"`
	TopicInputVolume3     mqtt.Topic    `mapstructure:"topic_volume_3" yaml:"topic_volume_3"`
	TopicInputVolume4     mqtt.Topic    `mapstructure:"topic_volume_4" yaml:"topic_volume_4"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/"

	return &Config{
		Location:              time.Now().Location().String(),
		UpdaterInterval:       time.Minute,
		TopicTemperatureIn:    prefix + "temperature_in",
		TopicTemperatureOut:   prefix + "temperature_out",
		TopicTemperatureDelta: prefix + "temperature_delta",
		TopicEnergy:           prefix + "energy",
		TopicConsumption:      prefix + "consumption",
		TopicCapacity:         prefix + "capacity",
		TopicPower:            prefix + "power",
		TopicInputPulses1:     prefix + "input/1/pulses",
		TopicInputPulses2:     prefix + "input/2/pulses",
		TopicInputPulses3:     prefix + "input/3/pulses",
		TopicInputPulses4:     prefix + "input/4/pulses",
		TopicInputVolume1:     prefix + "input/1/volume",
		TopicInputVolume2:     prefix + "input/2/volume",
		TopicInputVolume3:     prefix + "input/3/volume",
		TopicInputVolume4:     prefix + "input/4/volume",
	}
}

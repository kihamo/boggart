package pulsar

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	ConnectionDSN         string `mapstructure:"connection_dsn" yaml:"connection_dsn" valid:"required"`
	Address               string
	Location              string
	InputsCount           int64         `mapstructure:"inputs_count" yaml:"inputs_count"`
	StartPulsesInput1     int64         `mapstructure:"start_pulses_input_1" yaml:"start_pulses_input_1" valid:"float"`
	StartPulsesInput2     int64         `mapstructure:"start_pulses_input_2" yaml:"start_pulses_input_2" valid:"float"`
	StartPulsesInput3     int64         `mapstructure:"start_pulses_input_3" yaml:"start_pulses_input_3" valid:"float"`
	StartPulsesInput4     int64         `mapstructure:"start_pulses_input_4" yaml:"start_pulses_input_4" valid:"float"`
	StartVolumeInput1     float32       `mapstructure:"start_volume_input_1" yaml:"start_volume_input_1" valid:"float"`
	StartVolumeInput2     float32       `mapstructure:"start_volume_input_2" yaml:"start_volume_input_2" valid:"float"`
	StartVolumeInput3     float32       `mapstructure:"start_volume_input_3" yaml:"start_volume_input_3" valid:"float"`
	StartVolumeInput4     float32       `mapstructure:"start_volume_input_4" yaml:"start_volume_input_4" valid:"float"`
	UpdaterInterval       time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicTemperatureIn    mqtt.Topic    `mapstructure:"topic_temperature_in" yaml:"topic_temperature_in"`
	TopicTemperatureOut   mqtt.Topic    `mapstructure:"topic_temperature_out" yaml:"topic_temperature_out"`
	TopicTemperatureDelta mqtt.Topic    `mapstructure:"topic_temperature_delta" yaml:"topic_temperature_delta"`
	TopicEnergy           mqtt.Topic    `mapstructure:"topic_energy" yaml:"topic_energy"`
	TopicConsumption      mqtt.Topic    `mapstructure:"topic_consumption" yaml:"topic_consumption"`
	TopicCapacity         mqtt.Topic    `mapstructure:"topic_capacity" yaml:"topic_capacity"`
	TopicPower            mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicInputPulses      mqtt.Topic    `mapstructure:"topic_pulses" yaml:"topic_pulses"`
	TopicInputVolume      mqtt.Topic    `mapstructure:"topic_volume" yaml:"topic_volume"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.LivenessPeriod = time.Minute
	probesConfig.LivenessTimeout = time.Second * 10
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig:          probesConfig,
		LoggerConfig:          di.LoggerConfigDefaults(),
		InputsCount:           2,
		Location:              time.Now().Location().String(),
		UpdaterInterval:       time.Minute,
		TopicTemperatureIn:    prefix + "temperature_in",
		TopicTemperatureOut:   prefix + "temperature_out",
		TopicTemperatureDelta: prefix + "temperature_delta",
		TopicEnergy:           prefix + "energy",
		TopicConsumption:      prefix + "consumption",
		TopicCapacity:         prefix + "capacity",
		TopicPower:            prefix + "power",
		TopicInputPulses:      prefix + "input/pulses/+",
		TopicInputVolume:      prefix + "input/volume/+",
	}
}

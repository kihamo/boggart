package rpi

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	UpdaterInterval                        time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicModel                             mqtt.Topic    `mapstructure:"topic_model" yaml:"topic_model"`
	TopicCPUFrequentie                     mqtt.Topic    `mapstructure:"topic_cpu_frequentie" yaml:"topic_cpu_frequentie"`
	TopicTemperature                       mqtt.Topic    `mapstructure:"topic_temperature" yaml:"topic_temperature"`
	TopicVoltage                           mqtt.Topic    `mapstructure:"topic_voltage" yaml:"topic_voltage"`
	TopicCurrentlyUnderVoltage             mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicCurrentlyThrottled                mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicCurrentlyARMFrequencyCapped       mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicCurrentlySoftTemperatureReached   mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicSinceRebootUnderVoltage           mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicSinceRebootThrottled              mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicSinceRebootARMFrequencyCapped     mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
	TopicSinceRebootSoftTemperatureReached mqtt.Topic    `mapstructure:"topic_throttled" yaml:"topic_throttled"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/rpi/+/"

	return &Config{
		UpdaterInterval:                        time.Minute,
		TopicModel:                             prefix + "model",
		TopicCPUFrequentie:                     prefix + "cpu/+",
		TopicTemperature:                       prefix + "temperature",
		TopicVoltage:                           prefix + "voltage/+",
		TopicCurrentlyUnderVoltage:             prefix + "under-voltage/currently",
		TopicSinceRebootUnderVoltage:           prefix + "under-voltage/reboot",
		TopicCurrentlyThrottled:                prefix + "throttled/currently",
		TopicSinceRebootThrottled:              prefix + "throttled/reboot",
		TopicCurrentlyARMFrequencyCapped:       prefix + "arm-frequency-capped/currently",
		TopicSinceRebootARMFrequencyCapped:     prefix + "arm-frequency-capped/reboot",
		TopicCurrentlySoftTemperatureReached:   prefix + "soft-temperature-reached/currently",
		TopicSinceRebootSoftTemperatureReached: prefix + "soft-temperature-reached/reboot",
	}
}

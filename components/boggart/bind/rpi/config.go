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
	TopicCurrentlyUnderVoltage             mqtt.Topic    `mapstructure:"topic_currently_under_voltage" yaml:"topic_currently_under_voltage"`
	TopicCurrentlyThrottled                mqtt.Topic    `mapstructure:"topic_currently_throttled" yaml:"topic_currently_throttled"`
	TopicCurrentlyARMFrequencyCapped       mqtt.Topic    `mapstructure:"topic_currently_arm_frequency_capped" yaml:"topic_currently_arm_frequency_capped"`
	TopicCurrentlySoftTemperatureReached   mqtt.Topic    `mapstructure:"topic_currently_soft_temperature_reached" yaml:"topic_currently_soft_temperature_reached"`
	TopicSinceRebootUnderVoltage           mqtt.Topic    `mapstructure:"topic_since_reboot_under_voltage" yaml:"topic_since_reboot_under_voltage"`
	TopicSinceRebootThrottled              mqtt.Topic    `mapstructure:"topic_since_reboot_throttled" yaml:"topic_since_reboot_throttled"`
	TopicSinceRebootARMFrequencyCapped     mqtt.Topic    `mapstructure:"topic_since_reboot_arm_frequency_capped" yaml:"topic_since_reboot_arm_frequency_capped"`
	TopicSinceRebootSoftTemperatureReached mqtt.Topic    `mapstructure:"topic_since_reboot_soft_temperature_reached" yaml:"topic_since_reboot_soft_temperature_reached"`
}

func (Type) ConfigDefaults() interface{} {
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

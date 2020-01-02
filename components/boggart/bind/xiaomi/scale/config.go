package scale

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	MAC                  boggart.HardwareAddr `valid:",required"`
	UpdaterInterval      time.Duration        `mapstructure:"updater_interval" yaml:"updater_interval"`
	CaptureDuration      time.Duration        `mapstructure:"capture_interval" yaml:"capture_interval"`
	TopicWeight          mqtt.Topic           `mapstructure:"topic_weight" yaml:"topic_weight"`
	TopicImpedance       mqtt.Topic           `mapstructure:"topic_impedance" yaml:"topic_impedance"`
	TopicProfile         mqtt.Topic           `mapstructure:"topic_profile" yaml:"topic_profile"`
	TopicBMR             mqtt.Topic           `mapstructure:"topic_bmr" yaml:"topic_bmr"`
	TopicBMI             mqtt.Topic           `mapstructure:"topic_bmi" yaml:"topic_bmi"`
	TopicFatPercentage   mqtt.Topic           `mapstructure:"topic_fat_percentage" yaml:"topic_fat_percentage"`
	TopicWaterPercentage mqtt.Topic           `mapstructure:"topic_water_percentage" yaml:"topic_water_percentage"`
	TopicIdealWeight     mqtt.Topic           `mapstructure:"topic_ideal_weight" yaml:"topic_ideal_weight"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/xiaomi/scale/+/"

	return &Config{
		UpdaterInterval:      time.Minute,
		CaptureDuration:      time.Second * 10,
		TopicWeight:          prefix + "weight",
		TopicImpedance:       prefix + "impedance",
		TopicProfile:         prefix + "profile",
		TopicBMR:             prefix + "bmr",
		TopicBMI:             prefix + "bmi",
		TopicFatPercentage:   prefix + "fat-percentage",
		TopicWaterPercentage: prefix + "water-percentage",
		TopicIdealWeight:     prefix + "ideal-weight",
	}
}

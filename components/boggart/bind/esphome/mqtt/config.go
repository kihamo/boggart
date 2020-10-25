package mqtt

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	TopicDiscoveryPrefix  mqtt.Topic    `mapstructure:"topic_discovery_prefix" yaml:"topic_discovery_prefix"`
	TopicPrefix           mqtt.Topic    `valid:"required" mapstructure:"topic_prefix" yaml:"topic_prefix"`
	TopicLog              mqtt.Topic    `mapstructure:"topic_log" yaml:"topic_log"`
	TopicBirth            mqtt.Topic    `mapstructure:"topic_birth" yaml:"topic_birth"`
	TopicWill             mqtt.Topic    `mapstructure:"topic_will" yaml:"topic_will"`
	BirthMessage          string        `mapstructure:"birth_message" yaml:"birth_message"`
	WillMessage           string        `mapstructure:"will_message" yaml:"will_message"`
	IPAddressSensorID     string        `mapstructure:"ip_address_sensor_id" yaml:"ip_address_sensor_id"`
	ImportMetricsInterval time.Duration `mapstructure:"import_metrics_interval" yaml:"import_metrics_interval"`
}

func (t Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 15
	probesConfig.ReadinessTimeout = time.Second

	return &Config{
		ProbesConfig:          probesConfig,
		LoggerConfig:          di.LoggerConfigDefaults(),
		TopicDiscoveryPrefix:  "homeassistant",
		BirthMessage:          "online",
		WillMessage:           "offline",
		ImportMetricsInterval: time.Minute,
	}
}

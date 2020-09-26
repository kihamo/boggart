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
	TopicIPAddressSensor  mqtt.Topic    `mapstructure:"topic_ip_address_sensor" yaml:"topic_ip_address_sensor"`
	BirthMessage          string        `mapstructure:"birth_message" yaml:"birth_message"`
	WillMessage           string        `mapstructure:"will_message" yaml:"will_message"`
	ImportMetricsInterval time.Duration `mapstructure:"import_metrics_interval" yaml:"import_metrics_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 15,
			ReadinessTimeout: time.Second,
		},
		LoggerConfig:          di.LoggerConfigDefaults(),
		TopicDiscoveryPrefix:  "homeassistant",
		BirthMessage:          "online",
		WillMessage:           "offline",
		ImportMetricsInterval: time.Minute,
	}
}

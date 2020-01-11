package mqtt

import (
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	TopicDiscoveryPrefix mqtt.Topic `mapstructure:"topic_discovery_prefix" yaml:"topic_discovery_prefix"`
	TopicPrefix          mqtt.Topic `valid:"required" mapstructure:"topic_prefix" yaml:"topic_prefix"`
	TopicLog             mqtt.Topic `mapstructure:"topic_log" yaml:"topic_log"`
}

func (t Type) Config() interface{} {
	return &Config{
		TopicDiscoveryPrefix: "homeassistant",
	}
}

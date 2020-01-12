package mqtt

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	TopicDiscoveryPrefix mqtt.Topic `mapstructure:"topic_discovery_prefix" yaml:"topic_discovery_prefix"`
	TopicPrefix          mqtt.Topic `valid:"required" mapstructure:"topic_prefix" yaml:"topic_prefix"`
	TopicLog             mqtt.Topic `mapstructure:"topic_log" yaml:"topic_log"`
	TopicBirth           mqtt.Topic `mapstructure:"topic_birth" yaml:"topic_birth"`
	TopicWill            mqtt.Topic `mapstructure:"topic_will" yaml:"topic_will"`
	BirthMessage         string     `mapstructure:"birth_message" yaml:"birth_message"`
	WillMessage          string     `mapstructure:"will_message" yaml:"will_message"`
}

func (t Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Second * 15,
			ReadinessTimeout: time.Second,
		},
		TopicDiscoveryPrefix: "homeassistant",
		BirthMessage:         "online",
		WillMessage:          "offline",
	}
}

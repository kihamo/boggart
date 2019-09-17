package esphome

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Debug            bool
	Address          string `valid:"required"`
	Password         string
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicState       mqtt.Topic    `mapstructure:"topic_state" yaml:"topic_state"`
	TopicStateSet    mqtt.Topic    `mapstructure:"topic_state_set" yaml:"topic_state_set"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/esphome/+/"

	return &Config{
		Debug:            false,
		LivenessInterval: time.Minute,
		LivenessTimeout:  time.Second * 5,
		UpdaterInterval:  time.Minute,
		TopicState:       prefix + "+/state",
		TopicStateSet:    prefix + "+/set/state",
	}
}

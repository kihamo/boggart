package ds18b20

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address          string        `valid:"required"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicValue       mqtt.Topic    `mapstructure:"topic_value" yaml:"topic_value"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: time.Minute,
		LivenessTimeout:  time.Second * 5,
		UpdaterInterval:  time.Minute,
		TopicValue:       boggart.ComponentName + "/meter/ds18b20/+",
	}
}

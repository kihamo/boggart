package ds18b20

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Address         string        `valid:"required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicValue      mqtt.Topic    `mapstructure:"topic_value" yaml:"topic_value"`
}

func (t Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		UpdaterInterval: time.Minute,
		TopicValue:      boggart.ComponentName + "/meter/ds18b20/+",
	}
}

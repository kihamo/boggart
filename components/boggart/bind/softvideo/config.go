package softvideo

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`

	Login           string        `valid:"required"`
	Password        string        `valid:"required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	Debug           bool
	TopicBalance    mqtt.Topic `mapstructure:"topic_balance" yaml:"topic_balance"`
}

func (Type) Config() interface{} {
	return &Config{
		UpdaterInterval: time.Hour,
		Debug:           false,
		TopicBalance:    boggart.ComponentName + "/service/softvideo/+/balance",
	}
}

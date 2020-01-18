package rkcm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug           bool
	Login           string        `valid:"required"`
	Password        string        `valid:"required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicBalance    mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicMeterValue mqtt.Topic    `mapstructure:"topic_meter_value" yaml:"topic_meter_value"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/rkcm/+/"

	return &Config{
		Debug:           false,
		UpdaterInterval: time.Hour,
		TopicBalance:    prefix + "balance",
		TopicMeterValue: prefix + "meter/+",
	}
}

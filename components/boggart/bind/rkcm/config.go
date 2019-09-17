package rkcm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Debug           bool
	Login           string        `valid:"required"`
	Password        string        `valid:"required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicBalance    mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicMeterValue mqtt.Topic    `mapstructure:"topic_meter_value" yaml:"topic_meter_value"`
}

func (Type) Config() interface{} {
	return &Config{
		Debug:           false,
		UpdaterInterval: time.Hour,
		TopicBalance:    boggart.ComponentName + "/service/rkcm/+/balance",
		TopicMeterValue: boggart.ComponentName + "/service/rkcm/+/meter/+",
	}
}

package elektroset

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Debug               bool
	Login               string        `valid:"required"`
	Password            string        `valid:"required"`
	UpdaterInterval     time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicBalance        mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicServiceBalance mqtt.Topic    `mapstructure:"topic_service_balance" yaml:"topic_service_balance"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/elektroset/+/"

	return &Config{
		UpdaterInterval:     time.Hour,
		TopicBalance:        prefix + "balance",
		TopicServiceBalance: prefix + "+/balance",
	}
}

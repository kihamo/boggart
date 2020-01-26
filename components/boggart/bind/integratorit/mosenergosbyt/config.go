package mosenergosbyt

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug               bool
	Login               string        `valid:"required"`
	Password            string        `valid:"required"`
	UpdaterInterval     time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicBalance        mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicServiceBalance mqtt.Topic    `mapstructure:"topic_service_balance" yaml:"topic_service_balance"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/mosenergosbyt/+/"

	return &Config{
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		UpdaterInterval:     time.Hour,
		TopicBalance:        prefix + "balance",
		TopicServiceBalance: prefix + "balance/+",
	}
}

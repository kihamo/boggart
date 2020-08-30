package smtp

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	DSN       types.URL  `valid:"required"`
	Sender    string     `valid:"email,required"`
	TopicSend mqtt.Topic `mapstructure:"topic_send" yaml:"topic_send"`
}

func (t Type) Config() interface{} {
	return &Config{
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		TopicSend: boggart.ComponentName + "/mail/#",
	}
}

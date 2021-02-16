package smtp

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	DSN            types.URL  `valid:"required"`
	Sender         string     `valid:"email,required"`
	TopicSend      mqtt.Topic `mapstructure:"topic_send" yaml:"topic_send"`
	TopicSendMulti mqtt.Topic `mapstructure:"topic_send_multi" yaml:"topic_send_multi"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/mail/+/"

	return &Config{
		LoggerConfig:   di.LoggerConfigDefaults(),
		TopicSend:      prefix + "send",
		TopicSendMulti: prefix + "send/#",
	}
}

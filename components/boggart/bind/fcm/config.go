package fcm

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Tokens      []string   `valid:"required"`
	Credentials string     `valid:"required"`
	TopicSend   mqtt.Topic `mapstructure:"topic_send" yaml:"topic_send"`
}

func (Type) Config() interface{} {
	return &Config{
		TopicSend: boggart.ComponentName + "/fcm/+/send",
	}
}

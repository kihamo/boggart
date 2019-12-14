package fcm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Tokens      []string   `valid:"required"`
	Credentials string     `valid:"required"`
	TopicSend   mqtt.Topic `mapstructure:"topic_send" yaml:"topic_send"`
}

func (Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		TopicSend: boggart.ComponentName + "/fcm/+/send",
	}
}

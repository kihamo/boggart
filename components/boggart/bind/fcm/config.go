package fcm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Tokens      []string   `valid:"required"`
	Credentials string     `valid:"required"`
	TopicSend   mqtt.Topic `mapstructure:"topic_send" yaml:"topic_send"`
}

func (Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		TopicSend:    boggart.ComponentName + "/fcm/+/send",
	}
}

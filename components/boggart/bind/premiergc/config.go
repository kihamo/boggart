package premiergc

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Login        string `valid:"required"`
	Password     string `valid:"required"`
	Debug        bool
	TopicBalance mqtt.Topic `mapstructure:"topic_balance" yaml:"topic_balance"`
}

func (Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Hour
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		Debug:        false,
		TopicBalance: boggart.ComponentName + "/service/premiergc/+/balance",
	}
}

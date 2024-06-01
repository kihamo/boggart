package mosoblgaz

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Token        string     `valid:"required"`
	TopicBalance mqtt.Topic `mapstructure:"topic_balance" yaml:"topic_balance"`
}

func (Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/mosoblgaz/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Hour
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		TopicBalance: prefix + "balance",
	}
}

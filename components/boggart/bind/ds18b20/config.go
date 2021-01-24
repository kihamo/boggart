package ds18b20

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Sensors         []string
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicValue      mqtt.Topic    `mapstructure:"topic_value" yaml:"topic_value"`
}

func (t Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:    probesConfig,
		LoggerConfig:    di.LoggerConfigDefaults(),
		UpdaterInterval: time.Minute,
		TopicValue:      boggart.ComponentName + "/meter/ds18b20/+",
	}
}

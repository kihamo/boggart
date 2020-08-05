package influxdb

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	DSN             types.URL `valid:",required"`
	AuthToken       string    `mapstructure:"auth_token" yaml:"auth_token"`
	Organization    string
	Query           string        `valid:",required"`
	ExecuteInterval time.Duration `mapstructure:"execute_interval" yaml:"execute_interval"`
	ExecuteTimeout  time.Duration `mapstructure:"execute_timeout" yaml:"execute_timeout"`
	TopicResult     mqtt.Topic    `mapstructure:"topic_result" yaml:"topic_result"`
}

func (t Type) Config() interface{} {
	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		ExecuteInterval: time.Minute,
		ExecuteTimeout:  time.Second * 30,
		TopicResult:     boggart.ComponentName + "/influxdb/+",
	}
}

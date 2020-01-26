package nut

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address          boggart.URL   `valid:",required"`
	UPS              string        `valid:"required"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicVariable    mqtt.Topic    `mapstructure:"topic_variable" yaml:"topic_variable"`
	TopicVariableSet mqtt.Topic    `mapstructure:"topic_variable_set" yaml:"topic_variable_set"`
	TopicCommand     mqtt.Topic    `mapstructure:"topic_command" yaml:"topic_command"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/ups/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		UpdaterInterval:  time.Minute,
		TopicVariable:    prefix + "variable/+",
		TopicVariableSet: prefix + "variable/+/set",
		TopicCommand:     prefix + "command",
	}
}

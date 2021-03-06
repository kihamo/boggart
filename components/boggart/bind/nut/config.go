package nut

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

	Address          types.URL     `valid:",required"`
	UPS              string        `valid:"required"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicVariable    mqtt.Topic    `mapstructure:"topic_variable" yaml:"topic_variable"`
	TopicVariableSet mqtt.Topic    `mapstructure:"topic_variable_set" yaml:"topic_variable_set"`
	TopicCommand     mqtt.Topic    `mapstructure:"topic_command" yaml:"topic_command"`
	TopicCommandRun  mqtt.Topic    `mapstructure:"topic_command_run" yaml:"topic_command_run"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/ups/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:     probesConfig,
		LoggerConfig:     di.LoggerConfigDefaults(),
		UpdaterInterval:  time.Minute,
		TopicVariable:    prefix + "variable/+",
		TopicVariableSet: prefix + "variable/+/set",
		TopicCommand:     prefix + "command/+",
		TopicCommandRun:  prefix + "command/+/run",
	}
}

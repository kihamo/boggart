package myheat

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

	Address          types.URL `valid:",required"`
	Debug            bool
	TopicSensorValue mqtt.Topic    `mapstructure:"topic_sensor_value" yaml:"topic_sensor_value"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout   time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/myheat/+/"

	return &Config{
		ProbesConfig:     di.ProbesConfigDefaults(),
		LoggerConfig:     di.LoggerConfigDefaults(),
		UpdaterInterval:  time.Minute,
		TopicSensorValue: prefix + "sensor/+/value",
	}
}

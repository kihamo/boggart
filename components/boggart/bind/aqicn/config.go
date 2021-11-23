package aqicn

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug                   bool
	UpdaterInterval         time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	Token                   string        `mapstructure:"token" yaml:"token" valid:",required"`
	Latitude                float64
	Longitude               float64
	TopicCurrentTemperature mqtt.Topic `mapstructure:"topic_current_temperature" yaml:"topic_current_temperature"`
	TopicCurrentPressure    mqtt.Topic `mapstructure:"topic_current_pressure" yaml:"topic_current_pressure"`
	TopicCurrentHumidity    mqtt.Topic `mapstructure:"topic_current_humidity" yaml:"topic_current_humidity"`
	TopicCurrentDewPoint    mqtt.Topic `mapstructure:"topic_current_dew_point" yaml:"topic_current_dew_point"`
	TopicCurrentWindSpeed   mqtt.Topic `mapstructure:"topic_current_wind_speed" yaml:"topic_current_wind_speed"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/aqicn/+/"

	return &Config{
		LoggerConfig:            di.LoggerConfigDefaults(),
		UpdaterInterval:         time.Minute * 15,
		TopicCurrentTemperature: prefix + "current/temperature",
		TopicCurrentPressure:    prefix + "current/pressure",
		TopicCurrentHumidity:    prefix + "current/humidity",
		TopicCurrentDewPoint:    prefix + "current/dew-point",
		TopicCurrentWindSpeed:   prefix + "current/wind/speed",
	}
}

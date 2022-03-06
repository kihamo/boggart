package octoprint

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

	Address                types.URL `valid:",required"`
	APIKey                 string    `valid:",required" mapstructure:"api_key" yaml:"api_key"`
	Debug                  bool
	UpdaterInterval        time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout         time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	TopicState             mqtt.Topic    `mapstructure:"topic_state" yaml:"topic_state"`
	TopicTemperatureActual mqtt.Topic    `mapstructure:"topic_temperature_actual" yaml:"topic_temperature_actual"`
	TopicTemperatureOffset mqtt.Topic    `mapstructure:"topic_temperature_offset" yaml:"topic_temperature_offset"`
	TopicTemperatureTarget mqtt.Topic    `mapstructure:"topic_temperature_target" yaml:"topic_temperature_target"`
	TopicJobFileName       mqtt.Topic    `mapstructure:"topic_job_file_name" yaml:"topic_job_file_name"`
	TopicJobFileSize       mqtt.Topic    `mapstructure:"topic_job_file_size" yaml:"topic_job_file_size"`
	TopicJobProgress       mqtt.Topic    `mapstructure:"topic_job_progress" yaml:"topic_job_progress"`
	TopicJobTime           mqtt.Topic    `mapstructure:"topic_job_time" yaml:"topic_job_time"`
	TopicJobTimeLeft       mqtt.Topic    `mapstructure:"topic_job_time_left" yaml:"topic_job_time_left"`
	TopicLayerTotal        mqtt.Topic    `mapstructure:"topic_layer_total" yaml:"topic_layer_total"`
	TopicLayerCurrent      mqtt.Topic    `mapstructure:"topic_layer_current" yaml:"topic_layer_current"`
	TopicHeightTotal       mqtt.Topic    `mapstructure:"topic_height_total" yaml:"topic_height_total"`
	TopicHeightCurrent     mqtt.Topic    `mapstructure:"topic_height_current" yaml:"topic_height_current"`
	TopicCommand           mqtt.Topic    `mapstructure:"topic_command" yaml:"topic_command"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/octoprint/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 30
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:           probesConfig,
		LoggerConfig:           di.LoggerConfigDefaults(),
		UpdaterInterval:        time.Second * 30,
		UpdaterTimeout:         time.Second * 5,
		TopicState:             prefix + "state",
		TopicTemperatureActual: prefix + "temperature/+/actual",
		TopicTemperatureOffset: prefix + "temperature/+/offset",
		TopicTemperatureTarget: prefix + "temperature/+/target",
		TopicJobFileName:       prefix + "job/file/name",
		TopicJobFileSize:       prefix + "job/file/size",
		TopicJobProgress:       prefix + "job/progress",
		TopicJobTime:           prefix + "job/time",
		TopicJobTimeLeft:       prefix + "job/time-left",
		TopicLayerTotal:        prefix + "layer/total",
		TopicLayerCurrent:      prefix + "layer/current",
		TopicHeightTotal:       prefix + "height/total",
		TopicHeightCurrent:     prefix + "height/current",
		TopicCommand:           prefix + "command/+/+",
	}
}

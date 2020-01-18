package octoprint

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address                          boggart.URL `valid:",required"`
	APIKey                           string      `valid:",required" mapstructure:"api_key" yaml:"api_key"`
	Debug                            bool
	UpdaterInterval                  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout                   time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	TopicState                       mqtt.Topic    `mapstructure:"topic_state" yaml:"topic_state"`
	TopicStateBedTemperatureActual   mqtt.Topic    `mapstructure:"topic_state_bed_temperature_actual" yaml:"topic_state_bed_temperature_actual"`
	TopicStateBedTemperatureOffset   mqtt.Topic    `mapstructure:"topic_state_bed_temperature_offset" yaml:"topic_state_bed_temperature_offset"`
	TopicStateBedTemperatureTarget   mqtt.Topic    `mapstructure:"topic_state_bed_temperature_target" yaml:"topic_state_bed_temperature_target"`
	TopicStateTool0TemperatureActual mqtt.Topic    `mapstructure:"topic_state_tool0_temperature_actual" yaml:"topic_state_tool0_temperature_actual"`
	TopicStateTool0TemperatureOffset mqtt.Topic    `mapstructure:"topic_state_tool0_temperature_offset" yaml:"topic_state_tool0_temperature_offset"`
	TopicStateTool0TemperatureTarget mqtt.Topic    `mapstructure:"topic_state_tool0_temperature_target" yaml:"topic_state_tool0_temperature_target"`
	TopicStateJobFileName            mqtt.Topic    `mapstructure:"topic_state_job_file_name" yaml:"topic_state_job_file_name"`
	TopicStateJobFileSize            mqtt.Topic    `mapstructure:"topic_state_job_file_size" yaml:"topic_state_job_file_size"`
	TopicStateJobProgress            mqtt.Topic    `mapstructure:"topic_state_job_progress" yaml:"topic_state_job_progress"`
	TopicStateJobTime                mqtt.Topic    `mapstructure:"topic_state_job_time" yaml:"topic_state_job_time"`
	TopicStateJobTimeLeft            mqtt.Topic    `mapstructure:"topic_state_job_time_left" yaml:"topic_state_job_time_left"`
	TopicLayerTotal                  mqtt.Topic    `mapstructure:"topic_layer_total" yaml:"topic_layer_total"`
	TopicLayerCurrent                mqtt.Topic    `mapstructure:"topic_layer_current" yaml:"topic_layer_current"`
	TopicHeightTotal                 mqtt.Topic    `mapstructure:"topic_height_total" yaml:"topic_height_total"`
	TopicHeightCurrent               mqtt.Topic    `mapstructure:"topic_height_current" yaml:"topic_height_current"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/octoprint/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		UpdaterInterval:                  time.Second * 30,
		UpdaterTimeout:                   time.Second * 5,
		TopicState:                       prefix + "state",
		TopicStateBedTemperatureActual:   prefix + "state/bed/temperature/actual",
		TopicStateBedTemperatureOffset:   prefix + "state/bed/temperature/offset",
		TopicStateBedTemperatureTarget:   prefix + "state/bed/temperature/target",
		TopicStateTool0TemperatureActual: prefix + "state/tool0/temperature/actual",
		TopicStateTool0TemperatureOffset: prefix + "state/tool0/temperature/offset",
		TopicStateTool0TemperatureTarget: prefix + "state/tool0/temperature/target",
		TopicStateJobFileName:            prefix + "state/job/file/name",
		TopicStateJobFileSize:            prefix + "state/job/file/size",
		TopicStateJobProgress:            prefix + "state/job/progress",
		TopicStateJobTime:                prefix + "state/job/time",
		TopicStateJobTimeLeft:            prefix + "state/job/time-left",
		TopicLayerTotal:                  prefix + "layer/total",
		TopicLayerCurrent:                prefix + "layer/current",
		TopicHeightTotal:                 prefix + "height/total",
		TopicHeightCurrent:               prefix + "height/current",
	}
}

package octoprint

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address                          boggart.URL `valid:",required"`
	APIKey                           string      `valid:",required" mapstructure:"api_key" yaml:"api_key"`
	Debug                            bool
	LivenessInterval                 time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout                  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
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
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/octoprint/+/"

	return &Config{
		LivenessInterval:                 time.Second * 30,
		LivenessTimeout:                  time.Second * 5,
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
	}
}
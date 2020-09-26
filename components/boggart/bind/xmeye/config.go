package xmeye

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

	Address                        types.URL     `valid:",required"`
	UpdaterInterval                time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout                 time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	AlarmStreamingEnabled          bool          `mapstructure:"alarm_streaming_enabled" yaml:"alarm_streaming_enabled,omitempty"`
	AlarmStreamingInterval         time.Duration `mapstructure:"alarm_streaming_interval" yaml:"alarm_streaming_interval"`
	TopicEvent                     mqtt.Topic    `mapstructure:"topic_event" yaml:"topic_event"`
	TopicStateModel                mqtt.Topic    `mapstructure:"topic_state_model" yaml:"topic_state_model"`
	TopicStateFirmwareVersion      mqtt.Topic    `mapstructure:"topic_state_firmware_release_version" yaml:"topic_state_firmware_release_version"`
	TopicStateFirmwareReleasedDate mqtt.Topic    `mapstructure:"topic_state_firmware_release_date" yaml:"topic_state_firmware_release_date"`
	TopicStateUpTime               mqtt.Topic    `mapstructure:"topic_state_up_time" yaml:"topic_state_up_time"`
	TopicStateHDDCapacity          mqtt.Topic    `mapstructure:"topic_state_hdd_capacity" yaml:"topic_state_hdd_capacity"`
	TopicStateHDDFree              mqtt.Topic    `mapstructure:"topic_state_hdd_free" yaml:"topic_state_hdd_free"`
	TopicStateHDDUsage             mqtt.Topic    `mapstructure:"topic_state_hdd_usage" yaml:"topic_state_hdd_usage"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/"

	return &Config{
		LoggerConfig:                   di.LoggerConfigDefaults(),
		UpdaterInterval:                time.Minute,
		UpdaterTimeout:                 time.Second * 30,
		AlarmStreamingInterval:         time.Second * 5,
		TopicEvent:                     prefix + "+",
		TopicStateModel:                prefix + "state/model",
		TopicStateFirmwareVersion:      prefix + "state/firmware/version",
		TopicStateFirmwareReleasedDate: prefix + "state/firmware/release-date",
		TopicStateUpTime:               prefix + "state/uptime",
		TopicStateHDDCapacity:          prefix + "state/hdd/+/capacity",
		TopicStateHDDFree:              prefix + "state/hdd/+/free",
		TopicStateHDDUsage:             prefix + "state/hdd/+/usage",
	}
}

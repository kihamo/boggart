package tvt

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

	Address                   types.URL `valid:",required"`
	Debug                     bool
	UpdaterInterval           time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout            time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	TopicStateModel           mqtt.Topic    `mapstructure:"topic_state_model" yaml:"topic_state_model"`
	TopicStateFirmwareVersion mqtt.Topic    `mapstructure:"topic_state_firmware_release_version" yaml:"topic_state_firmware_release_version"`
	TopicStateHDDCapacity     mqtt.Topic    `mapstructure:"topic_state_hdd_capacity" yaml:"topic_state_hdd_capacity"`
	TopicStateHDDFree         mqtt.Topic    `mapstructure:"topic_state_hdd_free" yaml:"topic_state_hdd_free"`
	TopicStateHDDUsage        mqtt.Topic    `mapstructure:"topic_state_hdd_usage" yaml:"topic_state_hdd_usage"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/"

	return &Config{
		LoggerConfig:              di.LoggerConfigDefaults(),
		UpdaterInterval:           time.Minute,
		UpdaterTimeout:            time.Second * 30,
		TopicStateModel:           prefix + "state/model",
		TopicStateFirmwareVersion: prefix + "state/firmware/version",
		TopicStateHDDCapacity:     prefix + "state/hdd/+/capacity",
		TopicStateHDDFree:         prefix + "state/hdd/+/free",
		TopicStateHDDUsage:        prefix + "state/hdd/+/usage",
	}
}

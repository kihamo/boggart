package webos

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Host                    string        `valid:"host,required"`
	Key                     string        `valid:"required"`
	UpdaterInterval         time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicApplication        mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicMute               mqtt.Topic    `mapstructure:"topic_mute" yaml:"topic_mute"`
	TopicVolume             mqtt.Topic    `mapstructure:"topic_volume" yaml:"topic_volume"`
	TopicVolumeUp           mqtt.Topic    `mapstructure:"topic_volume_up" yaml:"topic_volume_up"`
	TopicVolumeDown         mqtt.Topic    `mapstructure:"topic_volume_down" yaml:"topic_volume_down"`
	TopicToast              mqtt.Topic    `mapstructure:"topic_toast" yaml:"topic_toast"`
	TopicPower              mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicStateMute          mqtt.Topic    `mapstructure:"topic_state_mute" yaml:"topic_state_mute"`
	TopicStateVolume        mqtt.Topic    `mapstructure:"topic_state_volume" yaml:"topic_state_volume"`
	TopicStateApplication   mqtt.Topic    `mapstructure:"topic_state_application" yaml:"topic_state_application"`
	TopicStateChannelNumber mqtt.Topic    `mapstructure:"topic_state_channel_number" yaml:"topic_state_channel_number"`
	TopicStatePower         mqtt.Topic    `mapstructure:"topic_state_power" yaml:"topic_state_power"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/tv/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 10,
			LivenessPeriod:   time.Second * 30,
			LivenessTimeout:  time.Second * 10,
		},
		LoggerConfig:            di.LoggerConfigDefaults(),
		UpdaterInterval:         time.Minute,
		TopicApplication:        prefix + "application",
		TopicMute:               prefix + "mute",
		TopicVolume:             prefix + "volume",
		TopicVolumeUp:           prefix + "volume/up",
		TopicVolumeDown:         prefix + "volume/down",
		TopicToast:              prefix + "toast",
		TopicPower:              prefix + "power",
		TopicStateMute:          prefix + "state/mute",
		TopicStateVolume:        prefix + "state/volume",
		TopicStateApplication:   prefix + "state/application",
		TopicStateChannelNumber: prefix + "state/channel-number",
		TopicStatePower:         prefix + "state/power",
	}
}

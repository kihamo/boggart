package alsa

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Volume           int64 `valid:"range(0|100)"`
	Mute             bool
	WidgetFileURL    string     `mapstructure:"widget_file_url" yaml:"widget_file_url"`
	TopicVolume      mqtt.Topic `mapstructure:"topic_volume" yaml:"topic_volume"`
	TopicMute        mqtt.Topic `mapstructure:"topic_mute" yaml:"topic_mute"`
	TopicPause       mqtt.Topic `mapstructure:"topic_pause" yaml:"topic_pause"`
	TopicStop        mqtt.Topic `mapstructure:"topic_stop" yaml:"topic_stop"`
	TopicPlay        mqtt.Topic `mapstructure:"topic_play" yaml:"topic_play"`
	TopicResume      mqtt.Topic `mapstructure:"topic_resume" yaml:"topic_resume"`
	TopicAction      mqtt.Topic `mapstructure:"topic_action" yaml:"topic_action"`
	TopicStateStatus mqtt.Topic `mapstructure:"topic_state_status" yaml:"topic_state_status"`
	TopicStateVolume mqtt.Topic `mapstructure:"topic_state_volume" yaml:"topic_state_volume"`
	TopicStateMute   mqtt.Topic `mapstructure:"topic_state_mute" yaml:"topic_state_mute"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/alsa/+/"

	return &Config{
		Volume:           50,
		Mute:             false,
		TopicVolume:      prefix + "volume",
		TopicMute:        prefix + "mute",
		TopicPause:       prefix + "pause",
		TopicStop:        prefix + "stop",
		TopicPlay:        prefix + "play",
		TopicResume:      prefix + "resume",
		TopicAction:      prefix + "action",
		TopicStateStatus: prefix + "state/status",
		TopicStateVolume: prefix + "state/volume",
		TopicStateMute:   prefix + "state/mute",
	}
}

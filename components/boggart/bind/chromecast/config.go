package chromecast

import (
	"time"

	"github.com/barnybug/go-cast/log"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug             bool
	Host              boggart.IP `valid:",required"`
	Port              int        `valid:"port"`
	Name              string
	WidgetFileURL     string     `mapstructure:"widget_file_url" yaml:"widget_file_url"`
	TopicVolume       mqtt.Topic `mapstructure:"topic_volume" yaml:"topic_volume"`
	TopicMute         mqtt.Topic `mapstructure:"topic_mute" yaml:"topic_mute"`
	TopicPause        mqtt.Topic `mapstructure:"topic_pause" yaml:"topic_pause"`
	TopicStop         mqtt.Topic `mapstructure:"topic_stop" yaml:"topic_stop"`
	TopicPlay         mqtt.Topic `mapstructure:"topic_play" yaml:"topic_play"`
	TopicResume       mqtt.Topic `mapstructure:"topic_resume" yaml:"topic_resume"`
	TopicSeek         mqtt.Topic `mapstructure:"topic_seek" yaml:"topic_seek"`
	TopicAction       mqtt.Topic `mapstructure:"topic_action" yaml:"topic_action"`
	TopicStateStatus  mqtt.Topic `mapstructure:"topic_state_status" yaml:"topic_state_status"`
	TopicStateVolume  mqtt.Topic `mapstructure:"topic_state_volume" yaml:"topic_state_volume"`
	TopicStateMute    mqtt.Topic `mapstructure:"topic_state_mute" yaml:"topic_state_mute"`
	TopicStateContent mqtt.Topic `mapstructure:"topic_state_content" yaml:"topic_state_content"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/chromecast/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 5,
			LivenessPeriod:   time.Second * 30,
			LivenessTimeout:  time.Second * 5,
		},
		Debug:             log.Debug,
		Port:              8009,
		Name:              boggart.ComponentName,
		TopicVolume:       prefix + "volume",
		TopicMute:         prefix + "mute",
		TopicPause:        prefix + "pause",
		TopicStop:         prefix + "stop",
		TopicPlay:         prefix + "play",
		TopicResume:       prefix + "resume",
		TopicSeek:         prefix + "seek",
		TopicAction:       prefix + "action",
		TopicStateStatus:  prefix + "state/status",
		TopicStateVolume:  prefix + "state/volume",
		TopicStateMute:    prefix + "state/mute",
		TopicStateContent: prefix + "state/content",
	}
}

package pass24online

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Phone           string `valid:",required"`
	Password        string `valid:",required"`
	Debug           bool
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicFeedEvent  mqtt.Topic    `mapstructure:"topic_feed_event" yaml:"topic_feed_event"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/pass24online/+/"

	return &Config{
		LoggerConfig:    di.LoggerConfigDefaults(),
		UpdaterInterval: time.Minute,
		TopicFeedEvent:  prefix + "feed/event",
	}
}

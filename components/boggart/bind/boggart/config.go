package boggart

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	ApplicationName    string     `valid:",required" mapstructure:"application_name" yaml:"application_name"`
	ApplicationVersion string     `valid:",required" mapstructure:"application_version" yaml:"application_version"`
	ApplicationBuild   string     `valid:",required" mapstructure:"application_build" yaml:"application_build"`
	TopicName          mqtt.Topic `mapstructure:"topic_name" yaml:"topic_name"`
	TopicVersion       mqtt.Topic `mapstructure:"topic_version" yaml:"topic_version"`
	TopicBuild         mqtt.Topic `mapstructure:"topic_build" yaml:"topic_build"`
	TopicShutdown      mqtt.Topic `mapstructure:"topic_shutdown" yaml:"topic_shutdown"`
}

func (Type) Config() interface{} {
	return &Config{
		TopicName:     boggart.ComponentName + "/boggart/+/application/name",
		TopicVersion:  boggart.ComponentName + "/boggart/+/application/version",
		TopicBuild:    boggart.ComponentName + "/boggart/+/application/build",
		TopicShutdown: boggart.ComponentName + "/boggart/+/application/shutdown",
	}
}

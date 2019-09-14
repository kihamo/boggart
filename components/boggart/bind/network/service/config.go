package service

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Hostname        string `valid:"host,required"`
	Port            int    `valid:"port,required"`
	Retry           int
	Timeout         time.Duration
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicOnline     mqtt.Topic    `mapstructure:"topic_online" yaml:"topic_online"`
	TopicLatency    mqtt.Topic    `mapstructure:"topic_latency" yaml:"topic_latency"`
	TopicCheck      mqtt.Topic    `mapstructure:"topic_check" yaml:"topic_check"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/+/"

	return &Config{
		Retry:           1,
		Timeout:         time.Second * 5,
		UpdaterInterval: time.Minute,
		TopicOnline:     prefix + "online",
		TopicLatency:    prefix + "latency",
		TopicCheck:      prefix + "check",
	}
}

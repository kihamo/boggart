package syslog

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Hostname     string
	Port         int64
	Timeout      time.Duration
	Tags         []string
	Topic        string
	TopicMessage mqtt.Topic `mapstructure:"topic_message" yaml:"topic_message"`
}

func (Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute * 10

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		Hostname:     "127.0.0.1",
		Port:         514,
		Timeout:      0,
		TopicMessage: boggart.ComponentName + "/syslog/+/message",
	}
}

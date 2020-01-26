package ping

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Hostname     string `valid:"host,required"`
	Retry        int
	TopicOnline  mqtt.Topic `mapstructure:"topic_online" yaml:"topic_online"`
	TopicLatency mqtt.Topic `mapstructure:"topic_latency" yaml:"topic_latency"`
	TopicCheck   mqtt.Topic `mapstructure:"topic_check" yaml:"topic_check"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/ping/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		Retry:        1,
		TopicOnline:  prefix + "online",
		TopicLatency: prefix + "latency",
		TopicCheck:   prefix + "check",
	}
}

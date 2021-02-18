package service

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
	Port         int    `valid:"port,required"`
	Retry        int
	TopicLatency mqtt.Topic `mapstructure:"topic_latency" yaml:"topic_latency"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		Retry:        1,
		TopicLatency: prefix + "latency",
	}
}

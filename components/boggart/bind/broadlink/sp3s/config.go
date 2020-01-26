package sp3s

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Host              string               `valid:",required"`
	MAC               boggart.HardwareAddr `valid:",required"`
	Model             string               `valid:"in(sp3seu|sp3sus),required"`
	UpdaterInterval   time.Duration        `mapstructure:"updater_interval" yaml:"updater_interval"`
	ConnectionTimeout time.Duration        `mapstructure:"connection_timeout" yaml:"connection_timeout"`
	TopicState        mqtt.Topic           `mapstructure:"topic_state" yaml:"topic_state"`
	TopicPower        mqtt.Topic           `mapstructure:"topic_power" yaml:"topic_power"`
	TopicSet          mqtt.Topic           `mapstructure:"topic_set" yaml:"topic_set"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/socket/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		UpdaterInterval:   time.Second * 3, // as e-control app, refresh every 3 sec,
		ConnectionTimeout: time.Second,
		TopicState:        prefix + "state",
		TopicPower:        prefix + "power",
		TopicSet:          prefix + "set",
	}
}

package sp3s

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
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
		UpdaterInterval:   time.Second * 3, // as e-control app, refresh every 3 sec,
		ConnectionTimeout: time.Second,
		TopicState:        prefix + "state",
		TopicPower:        prefix + "power",
		TopicSet:          prefix + "set",
	}
}

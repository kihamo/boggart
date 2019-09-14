package gpio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Pin           uint64     `valid:"required"`
	Mode          string     `valid:"in(in|out)"`
	TopicPinState mqtt.Topic `mapstructure:"topic_pin_state" yaml:"topic_pin_state"`
	TopicPinSet   mqtt.Topic `mapstructure:"topic_pin_set" yaml:"topic_pin_set"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/gpio/+"

	return &Config{
		TopicPinState: prefix,
		TopicPinSet:   prefix + "/set",
	}
}

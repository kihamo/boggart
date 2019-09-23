package led_wifi

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address            string        `valid:"host,required"`
	UpdaterInterval    time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicPower         mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicColor         mqtt.Topic    `mapstructure:"topic_color" yaml:"topic_color"`
	TopicMode          mqtt.Topic    `mapstructure:"topic_mode" yaml:"topic_mode"`
	TopicSpeed         mqtt.Topic    `mapstructure:"topic_speed" yaml:"topic_speed"`
	TopicStatePower    mqtt.Topic    `mapstructure:"topic_state_power" yaml:"topic_state_power"`
	TopicStateColor    mqtt.Topic    `mapstructure:"topic_state_color" yaml:"topic_state_color"`
	TopicStateColorHSV mqtt.Topic    `mapstructure:"topic_state_color_hsv" yaml:"topic_state_color_hsv"`
	TopicStateMode     mqtt.Topic    `mapstructure:"topic_state_mode" yaml:"topic_state_mode"`
	TopicStateSpeed    mqtt.Topic    `mapstructure:"topic_state_speed" yaml:"topic_state_speed"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/led/+/"

	return &Config{
		UpdaterInterval:    time.Second * 3,
		TopicPower:         prefix + "power",
		TopicColor:         prefix + "color",
		TopicMode:          prefix + "mode",
		TopicSpeed:         prefix + "speed",
		TopicStatePower:    prefix + "state/power",
		TopicStateColor:    prefix + "state/color",
		TopicStateColorHSV: prefix + "state/color/hsv",
		TopicStateMode:     prefix + "state/mode",
		TopicStateSpeed:    prefix + "state/speed",
	}
}

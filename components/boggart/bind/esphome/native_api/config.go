package nativeapi

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug                      bool
	Address                    string `valid:"required"`
	Password                   string
	OTAPort                    uint64        `mapstructure:"ota_port" yaml:"ota_port"`
	OTAPassword                string        `mapstructure:"ota_password" yaml:"ota_password"`
	SyncStateInterval          time.Duration `mapstructure:"sync_state_interval" yaml:"sync_state_interval"`
	TopicPower                 mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicColor                 mqtt.Topic    `mapstructure:"topic_color" yaml:"topic_color"`
	TopicState                 mqtt.Topic    `mapstructure:"topic_state" yaml:"topic_state"`
	TopicStateColorRGB         mqtt.Topic    `mapstructure:"topic_state_color_rgb" yaml:"topic_state_color_rgb"`
	TopicStateColorHSV         mqtt.Topic    `mapstructure:"topic_state_color_hsv" yaml:"topic_state_color_hsv"`
	TopicStateBrightness       mqtt.Topic    `mapstructure:"topic_state_brightness" yaml:"topic_state_brightness"`
	TopicStateRed              mqtt.Topic    `mapstructure:"topic_state_red" yaml:"topic_state_red"`
	TopicStateGreen            mqtt.Topic    `mapstructure:"topic_state_green" yaml:"topic_state_green"`
	TopicStateBlue             mqtt.Topic    `mapstructure:"topic_state_blue" yaml:"topic_state_blue"`
	TopicStateWhite            mqtt.Topic    `mapstructure:"topic_state_white" yaml:"topic_state_white"`
	TopicStateColorTemperature mqtt.Topic    `mapstructure:"topic_state_color_temperature" yaml:"topic_state_color_temperature"`
	TopicStateEffect           mqtt.Topic    `mapstructure:"topic_state_effect" yaml:"topic_state_effect"`
	TopicStateSet              mqtt.Topic    `mapstructure:"topic_state_set" yaml:"topic_state_set"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/esphome/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:               probesConfig,
		LoggerConfig:               di.LoggerConfigDefaults(),
		Debug:                      false,
		SyncStateInterval:          time.Minute,
		TopicPower:                 prefix + "+/power",
		TopicColor:                 prefix + "+/color",
		TopicState:                 prefix + "+/state",
		TopicStateColorRGB:         prefix + "+/state/rgb",
		TopicStateColorHSV:         prefix + "+/state/hsv",
		TopicStateBrightness:       prefix + "+/state/brightness",
		TopicStateRed:              prefix + "+/state/red",
		TopicStateGreen:            prefix + "+/state/green",
		TopicStateBlue:             prefix + "+/state/blue",
		TopicStateWhite:            prefix + "+/state/white",
		TopicStateColorTemperature: prefix + "+/state/color-temperature",
		TopicStateEffect:           prefix + "+/state/effect",
		TopicStateSet:              prefix + "+/set/state",
	}
}

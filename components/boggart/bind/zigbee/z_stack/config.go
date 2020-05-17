package z_stack

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	ConnectionDSN       string     `mapstructure:"connection_dsn" yaml:"connection_dsn" valid:"required"`
	TopicLinkQuality    mqtt.Topic `mapstructure:"topic_link_quality" yaml:"topic_link_quality"`
	TopicBatteryPercent mqtt.Topic `mapstructure:"topic_battery_percent" yaml:"topic_battery_percent"`
	TopicBatteryVoltage mqtt.Topic `mapstructure:"topic_battery_voltage" yaml:"topic_battery_voltage"`
	TopicOnOff          mqtt.Topic `mapstructure:"topic_on_off" yaml:"topic_on_off"`
	TopicClick          mqtt.Topic `mapstructure:"topic_click" yaml:"topic_click"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/zigbee/zstack/+/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			LivenessPeriod:  time.Minute * 10,
			LivenessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		TopicLinkQuality:    prefix + "link-quality",
		TopicBatteryPercent: prefix + "battery/percent",
		TopicBatteryVoltage: prefix + "battery/voltage",
		TopicOnOff:          prefix + "on-off",
		TopicClick:          prefix + "click",
	}
}

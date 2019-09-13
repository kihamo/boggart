package miio

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Host                     string        `valid:"host,required"`
	Token                    string        `valid:"required"`
	PacketsCounter           uint32        `mapstructure:"packets_counter" yaml:"packets_counter"`
	LivenessInterval         time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout          time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval          time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout           time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	TopicBattery             mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicCleanArea           mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicCleanTime           mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicFanPower            mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicVolume              mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicConsumableFilter    mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicConsumableBrushMain mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicConsumableBrushSide mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicConsumableSensor    mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicState               mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicError               mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicSetFanPower         mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicSetVolume           mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicTestVolume          mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicFind                mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicAction              mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/xiaomi/roborock/+/"

	return &Config{
		LivenessInterval:         time.Minute,
		LivenessTimeout:          time.Second * 5,
		UpdaterInterval:          time.Minute,
		UpdaterTimeout:           time.Second * 30,
		TopicBattery:             prefix + "battery",
		TopicCleanArea:           prefix + "clean/area",
		TopicCleanTime:           prefix + "clean/time",
		TopicFanPower:            prefix + "fan-power",
		TopicVolume:              prefix + "volume",
		TopicConsumableFilter:    prefix + "consumable/filter",
		TopicConsumableBrushMain: prefix + "consumable/brush-main",
		TopicConsumableBrushSide: prefix + "consumable/brush-side",
		TopicConsumableSensor:    prefix + "consumable/sensor",
		TopicState:               prefix + "state",
		TopicError:               prefix + "error",
		TopicSetFanPower:         prefix + "fan-power/set",
		TopicSetVolume:           prefix + "volume/set",
		TopicTestVolume:          prefix + "volume/test",
		TopicFind:                prefix + "find",
		TopicAction:              prefix + "action",
	}
}

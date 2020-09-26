package miio

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Host                        string        `valid:"host,required"`
	Token                       string        `valid:"required"`
	PacketsCounter              uint32        `mapstructure:"packets_counter" yaml:"packets_counter"`
	UpdaterInterval             time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout              time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	TopicBattery                mqtt.Topic    `mapstructure:"topic_battery" yaml:"topic_battery"`
	TopicLastCleanCompleted     mqtt.Topic    `mapstructure:"topic_last_clean_completed" yaml:"topic_last_clean_completed"`
	TopicLastCleanArea          mqtt.Topic    `mapstructure:"topic_last_clean_area" yaml:"topic_last_clean_area"`
	TopicLastCleanStartDateTime mqtt.Topic    `mapstructure:"topic_last_clean_start_datetime" yaml:"topic_last_clean_start_datetime"`
	TopicLastCleanEndDateTime   mqtt.Topic    `mapstructure:"topic_last_clean_end_datetime" yaml:"topic_last_clean_end_datetime"`
	TopicLastCleanDuration      mqtt.Topic    `mapstructure:"topic_last_clean_duration" yaml:"topic_last_clean_duration"`
	TopicFanPower               mqtt.Topic    `mapstructure:"topic_fan_power" yaml:"topic_fan_power"`
	TopicVolume                 mqtt.Topic    `mapstructure:"topic_volume" yaml:"topic_volume"`
	TopicConsumableFilter       mqtt.Topic    `mapstructure:"topic_consumable_filter" yaml:"topic_consumable_filter"`
	TopicConsumableBrushMain    mqtt.Topic    `mapstructure:"topic_consumable_brush_main" yaml:"topic_consumable_brush_main"`
	TopicConsumableBrushSide    mqtt.Topic    `mapstructure:"topic_consumable_brush_side" yaml:"topic_consumable_brush_side"`
	TopicConsumableSensor       mqtt.Topic    `mapstructure:"topic_consumable_sensor" yaml:"topic_consumable_sensor"`
	TopicState                  mqtt.Topic    `mapstructure:"topic_state" yaml:"topic_state"`
	TopicError                  mqtt.Topic    `mapstructure:"topic_error" yaml:"topic_error"`
	TopicSetFanPower            mqtt.Topic    `mapstructure:"topic_set_fan_power" yaml:"topic_set_fan_power"`
	TopicSetVolume              mqtt.Topic    `mapstructure:"topic_set_volume" yaml:"topic_set_volume"`
	TopicTestVolume             mqtt.Topic    `mapstructure:"topic_test_volume" yaml:"topic_test_volume"`
	TopicFind                   mqtt.Topic    `mapstructure:"topic_find" yaml:"topic_find"`
	TopicAction                 mqtt.Topic    `mapstructure:"topic_action" yaml:"topic_action"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/xiaomi/roborock/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig:                di.LoggerConfigDefaults(),
		UpdaterInterval:             time.Minute,
		UpdaterTimeout:              time.Second * 30,
		TopicBattery:                prefix + "battery",
		TopicLastCleanCompleted:     prefix + "clean/last/completed",
		TopicLastCleanArea:          prefix + "clean/last/area",
		TopicLastCleanStartDateTime: prefix + "clean/last/start-time",
		TopicLastCleanEndDateTime:   prefix + "clean/last/end-time",
		TopicLastCleanDuration:      prefix + "clean/last/duration",
		TopicFanPower:               prefix + "fan-power",
		TopicVolume:                 prefix + "volume",
		TopicConsumableFilter:       prefix + "consumable/filter",
		TopicConsumableBrushMain:    prefix + "consumable/brush-main",
		TopicConsumableBrushSide:    prefix + "consumable/brush-side",
		TopicConsumableSensor:       prefix + "consumable/sensor",
		TopicState:                  prefix + "state",
		TopicError:                  prefix + "error",
		TopicSetFanPower:            prefix + "fan-power/set",
		TopicSetVolume:              prefix + "volume/set",
		TopicTestVolume:             prefix + "volume/test",
		TopicFind:                   prefix + "find",
		TopicAction:                 prefix + "action",
	}
}

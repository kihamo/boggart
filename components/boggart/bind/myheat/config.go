package myheat

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address                  types.URL `valid:",required"`
	Debug                    bool
	TopicAlarmPowerSupply    mqtt.Topic    `mapstructure:"topic_alarm_power_supply" yaml:"topic_alarm_power_supply"`
	TopicAlarmReplaceBattery mqtt.Topic    `mapstructure:"topic_alarm_replace_battery" yaml:"topic_alarm_replace_battery"`
	TopicAlarmGSMBalance     mqtt.Topic    `mapstructure:"topic_alarm_gsm_balance" yaml:"topic_alarm_gsm_balance"`
	TopicInternetConnected   mqtt.Topic    `mapstructure:"topic_internet_connected" yaml:"topic_internet_connected"`
	TopicDeviceSeverity      mqtt.Topic    `mapstructure:"topic_device_severity" yaml:"topic_device_severity"`
	TopicGSMSignalLevel      mqtt.Topic    `mapstructure:"topic_gsm_signal_level" yaml:"topic_gsm_signal_level"`
	TopicGSMBalance          mqtt.Topic    `mapstructure:"topic_gsm_balance" yaml:"topic_gsm_balance"`
	TopicSecurityArmed       mqtt.Topic    `mapstructure:"topic_security_armed" yaml:"topic_security_armed"`
	TopicSecurityArmedState  mqtt.Topic    `mapstructure:"topic_security_armed_state" yaml:"topic_security_armed_state"`
	TopicSensorValue         mqtt.Topic    `mapstructure:"topic_sensor_value" yaml:"topic_sensor_value"`
	UpdaterInterval          time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout           time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/myheat/+/"

	return &Config{
		ProbesConfig:             di.ProbesConfigDefaults(),
		LoggerConfig:             di.LoggerConfigDefaults(),
		UpdaterInterval:          time.Minute,
		TopicAlarmPowerSupply:    prefix + "alarm/power-supply",
		TopicAlarmReplaceBattery: prefix + "alarm/replace-battery",
		TopicAlarmGSMBalance:     prefix + "alarm/gsm-balance",
		TopicInternetConnected:   prefix + "device/internet",
		TopicDeviceSeverity:      prefix + "device/severity",
		TopicGSMSignalLevel:      prefix + "gsm/signal",
		TopicGSMBalance:          prefix + "gsm/balance",
		TopicSecurityArmed:       prefix + "security",
		TopicSecurityArmedState:  prefix + "security/state",
		TopicSensorValue:         prefix + "sensor/+/value",
	}
}

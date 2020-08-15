package v3

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	ConnectionDSN        string `mapstructure:"connection_dsn" yaml:"connection_dsn" valid:"required"`
	Address              string
	UpdaterInterval      time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicTariff1         mqtt.Topic    `mapstructure:"topic_tariff_1" yaml:"topic_tariff_1"`
	TopicVoltage1        mqtt.Topic    `mapstructure:"topic_voltage_1" yaml:"topic_voltage_1"`
	TopicVoltage2        mqtt.Topic    `mapstructure:"topic_voltage_2" yaml:"topic_voltage_2"`
	TopicVoltage3        mqtt.Topic    `mapstructure:"topic_voltage_3" yaml:"topic_voltage_3"`
	TopicAmperage1       mqtt.Topic    `mapstructure:"topic_amperage_1" yaml:"topic_amperage_1"`
	TopicAmperage2       mqtt.Topic    `mapstructure:"topic_amperage_2" yaml:"topic_amperage_2"`
	TopicAmperage3       mqtt.Topic    `mapstructure:"topic_amperage_3" yaml:"topic_amperage_3"`
	TopicPower1          mqtt.Topic    `mapstructure:"topic_power_1" yaml:"topic_power_1"`
	TopicPower2          mqtt.Topic    `mapstructure:"topic_power_2" yaml:"topic_power_2"`
	TopicPower3          mqtt.Topic    `mapstructure:"topic_power_3" yaml:"topic_power_3"`
	TopicMakeDate        mqtt.Topic    `mapstructure:"topic_make_date" yaml:"topic_make_date"`
	TopicFirmwareVersion mqtt.Topic    `mapstructure:"topic_firmware_version" yaml:"topic_firmware_version"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			LivenessPeriod:   time.Minute,
			LivenessTimeout:  time.Second * 10,
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 10,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		UpdaterInterval:      time.Minute * 5,
		TopicTariff1:         prefix + "tariff/1",
		TopicVoltage1:        prefix + "voltage/1",
		TopicVoltage2:        prefix + "voltage/2",
		TopicVoltage3:        prefix + "voltage/3",
		TopicAmperage1:       prefix + "amperage/1",
		TopicAmperage2:       prefix + "amperage/2",
		TopicAmperage3:       prefix + "amperage/3",
		TopicPower1:          prefix + "power/1",
		TopicPower2:          prefix + "power/2",
		TopicPower3:          prefix + "power/3",
		TopicMakeDate:        prefix + "make-date",
		TopicFirmwareVersion: prefix + "firmware/version",
	}
}

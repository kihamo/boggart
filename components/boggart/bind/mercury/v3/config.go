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
	TopicTariff          mqtt.Topic    `mapstructure:"topic_tariff" yaml:"topic_tariff"`
	TopicVoltagePhase1   mqtt.Topic    `mapstructure:"topic_voltage_phase_1" yaml:"topic_voltage_phase_1"`
	TopicVoltagePhase2   mqtt.Topic    `mapstructure:"topic_voltage_phase_2" yaml:"topic_voltage_phase_2"`
	TopicVoltagePhase3   mqtt.Topic    `mapstructure:"topic_voltage_phase_3" yaml:"topic_voltage_phase_3"`
	TopicAmperagePhase1  mqtt.Topic    `mapstructure:"topic_amperage_phase_1" yaml:"topic_amperage_phase_1"`
	TopicAmperagePhase2  mqtt.Topic    `mapstructure:"topic_amperage_phase_2" yaml:"topic_amperage_phase_2"`
	TopicAmperagePhase3  mqtt.Topic    `mapstructure:"topic_amperage_phase_3" yaml:"topic_amperage_phase_3"`
	TopicAmperageTotal   mqtt.Topic    `mapstructure:"topic_amperage_total" yaml:"topic_amperage_total"`
	TopicPowerPhase1     mqtt.Topic    `mapstructure:"topic_power_phase_1" yaml:"topic_power_phase_1"`
	TopicPowerPhase2     mqtt.Topic    `mapstructure:"topic_power_phase_2" yaml:"topic_power_phase_2"`
	TopicPowerPhase3     mqtt.Topic    `mapstructure:"topic_power_phase_3" yaml:"topic_power_phase_3"`
	TopicPowerTotal      mqtt.Topic    `mapstructure:"topic_power_total" yaml:"topic_power_total"`
	TopicMakeDate        mqtt.Topic    `mapstructure:"topic_make_date" yaml:"topic_make_date"`
	TopicFirmwareVersion mqtt.Topic    `mapstructure:"topic_firmware_version" yaml:"topic_firmware_version"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.LivenessPeriod = time.Minute
	probesConfig.LivenessTimeout = time.Second * 10
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig:         probesConfig,
		LoggerConfig:         di.LoggerConfigDefaults(),
		UpdaterInterval:      time.Minute * 5,
		TopicTariff:          prefix + "tariff/1",
		TopicVoltagePhase1:   prefix + "voltage/1",
		TopicVoltagePhase2:   prefix + "voltage/2",
		TopicVoltagePhase3:   prefix + "voltage/3",
		TopicAmperagePhase1:  prefix + "amperage/1",
		TopicAmperagePhase2:  prefix + "amperage/2",
		TopicAmperagePhase3:  prefix + "amperage/3",
		TopicAmperageTotal:   prefix + "amperage/total",
		TopicPowerPhase1:     prefix + "power/1",
		TopicPowerPhase2:     prefix + "power/2",
		TopicPowerPhase3:     prefix + "power/3",
		TopicPowerTotal:      prefix + "power/total",
		TopicMakeDate:        prefix + "make-date",
		TopicFirmwareVersion: prefix + "firmware/version",
	}
}

package v1

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
	Address              string `valid:"required"`
	Location             string
	UpdaterInterval      time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicTariff          mqtt.Topic    `mapstructure:"topic_tariff" yaml:"topic_tariff"`
	TopicVoltage         mqtt.Topic    `mapstructure:"topic_voltage" yaml:"topic_voltage"`
	TopicAmperage        mqtt.Topic    `mapstructure:"topic_amperage" yaml:"topic_amperage"`
	TopicPower           mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicBatteryVoltage  mqtt.Topic    `mapstructure:"topic_battery_voltage" yaml:"topic_battery_voltage"`
	TopicLastPowerOff    mqtt.Topic    `mapstructure:"topic_last_power_off" yaml:"topic_last_power_off"`
	TopicLastPowerOn     mqtt.Topic    `mapstructure:"topic_last_power_on" yaml:"topic_last_power_on"`
	TopicMakeDate        mqtt.Topic    `mapstructure:"topic_make_date" yaml:"topic_make_date"`
	TopicFirmwareDate    mqtt.Topic    `mapstructure:"topic_firmware_date" yaml:"topic_firmware_date"`
	TopicFirmwareVersion mqtt.Topic    `mapstructure:"topic_firmware_version" yaml:"topic_firmware_version"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/meter/mercury/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.LivenessPeriod = time.Minute
	probesConfig.LivenessTimeout = time.Second * 10
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		Location:     time.Now().Location().String(),
		/*
			При отсутствии тока в последовательной цепи и значении напряжения, равном 1,15Uном, испытательный выход
			счётчика не создаёт более одного импульса в течение времени, равного 4,4 мин и 3,5 мин для счётчиков класса
			точности 1 и 2 соответственно.
		*/
		UpdaterInterval:      time.Minute * 5,
		TopicTariff:          prefix + "tariff/+",
		TopicVoltage:         prefix + "voltage",
		TopicAmperage:        prefix + "amperage",
		TopicPower:           prefix + "power",
		TopicBatteryVoltage:  prefix + "battery_voltage",
		TopicLastPowerOff:    prefix + "last-power-off",
		TopicLastPowerOn:     prefix + "last-power-on",
		TopicMakeDate:        prefix + "make-date",
		TopicFirmwareDate:    prefix + "firmware/date",
		TopicFirmwareVersion: prefix + "firmware/version",
	}
}

package modbus

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

	DSN                   types.URL     `valid:"required"`
	ConnectionSlaveID     uint8         `mapstructure:"connection_slave_id" yaml:"connection_slave_id"`
	ConnectionTimeout     time.Duration `mapstructure:"connection_timeout" yaml:"connection_timeout"`
	ConnectionIdleTimeout time.Duration `mapstructure:"connection_idle_timeout" yaml:"connection_idle_timeout"`
	SensorUpdaterInterval time.Duration `mapstructure:"sensor_updater_interval" yaml:"sensor_updater_interval"`
	StatusUpdaterInterval time.Duration `mapstructure:"status_updater_interval" yaml:"status_updater_interval"`
	//DefaultsTemperatureFormat uint16        `mapstructure:"defaults_temperature_format" yaml:"defaults_temperature_format" valid:"range(0|1)"`
	DefaultsAwayTemperature      uint16     `mapstructure:"defaults_away_temperature" yaml:"defaults_away_temperature" valid:"range(5|35)"`
	TopicDeviceType              mqtt.Topic `mapstructure:"topic_device_type" yaml:"topic_device_type"`
	TopicHeatingOutputStatus     mqtt.Topic `mapstructure:"topic_heating_output_status" yaml:"topic_heating_output_status"`
	TopicHoldingFunction         mqtt.Topic `mapstructure:"topic_holding_function" yaml:"topic_holding_function"`
	TopicFloorOverheat           mqtt.Topic `mapstructure:"topic_floor_overheat" yaml:"topic_floor_overheat"`
	TopicRoomTemperature         mqtt.Topic `mapstructure:"topic_room_temperature" yaml:"topic_room_temperature"`
	TopicFloorTemperature        mqtt.Topic `mapstructure:"topic_floor_temperature" yaml:"topic_floor_temperature"`
	TopicHumidity                mqtt.Topic `mapstructure:"topic_humidity" yaml:"topic_humidity"`
	TopicPower                   mqtt.Topic `mapstructure:"topic_power" yaml:"topic_power"`
	TopicPowerState              mqtt.Topic `mapstructure:"topic_power_state" yaml:"topic_power_state"`
	TopicSetTemperature          mqtt.Topic `mapstructure:"topic_set_temperature" yaml:"topic_set_temperature"`
	TopicSetTemperatureState     mqtt.Topic `mapstructure:"topic_set_temperature_state" yaml:"topic_set_temperature_state"`
	TopicAway                    mqtt.Topic `mapstructure:"topic_away" yaml:"topic_away"`
	TopicAwayState               mqtt.Topic `mapstructure:"topic_away_state" yaml:"topic_away_state"`
	TopicAwayTemperature         mqtt.Topic `mapstructure:"topic_away_temperature" yaml:"topic_away_temperature"`
	TopicAwayTemperatureState    mqtt.Topic `mapstructure:"topic_away_temperature_state" yaml:"topic_away_temperature_state"`
	TopicHoldingTemperature      mqtt.Topic `mapstructure:"topic_holding_temperature" yaml:"topic_holding_temperature"`
	TopicHoldingTemperatureState mqtt.Topic `mapstructure:"topic_holding_temperature_state" yaml:"topic_holding_temperature_state"`
	//TopicTemperatureFormat      mqtt.Topic    `mapstructure:"topic_temperature_format" yaml:"topic_temperature_format"`
	//TopicTemperatureFormatState mqtt.Topic    `mapstructure:"topic_temperature_format_state" yaml:"topic_temperature_format_state"`
}

func (t Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	var prefix mqtt.Topic = boggart.ComponentName + "/mc6/+/"

	return &Config{
		ProbesConfig:          probesConfig,
		LoggerConfig:          di.LoggerConfigDefaults(),
		ConnectionSlaveID:     0x1,
		ConnectionTimeout:     time.Second,
		ConnectionIdleTimeout: time.Minute,
		SensorUpdaterInterval: time.Minute,
		StatusUpdaterInterval: time.Second * 10,
		//DefaultsTemperatureFormat:   mc6.TemperatureFormatCelsius,
		DefaultsAwayTemperature:      7,
		TopicDeviceType:              prefix + "type",
		TopicHeatingOutputStatus:     prefix + "heating/state",
		TopicHoldingFunction:         prefix + "holding/state",
		TopicFloorOverheat:           prefix + "floor/overheat",
		TopicRoomTemperature:         prefix + "room/temperature/state",
		TopicFloorTemperature:        prefix + "floor/temperature/state",
		TopicHumidity:                prefix + "humidity/state",
		TopicPower:                   prefix + "power",
		TopicPowerState:              prefix + "power/state",
		TopicSetTemperature:          prefix + "temperature",
		TopicSetTemperatureState:     prefix + "temperature/state",
		TopicAway:                    prefix + "away",
		TopicAwayState:               prefix + "away/state",
		TopicAwayTemperature:         prefix + "away/temperature",
		TopicAwayTemperatureState:    prefix + "away/temperature/state",
		TopicHoldingTemperature:      prefix + "holding/temperature",
		TopicHoldingTemperatureState: prefix + "holding/temperature/state",
		//TopicTemperatureFormat:      prefix + "temperature/format",
		//TopicTemperatureFormatState: prefix + "temperature/format/state",
	}
}

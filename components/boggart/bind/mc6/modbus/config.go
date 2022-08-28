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

	DSN                      types.URL     `valid:"required"`
	ConnectionSlaveID        uint8         `mapstructure:"connection_slave_id" yaml:"connection_slave_id"`
	ConnectionTimeout        time.Duration `mapstructure:"connection_timeout" yaml:"connection_timeout"`
	ConnectionIdleTimeout    time.Duration `mapstructure:"connection_idle_timeout" yaml:"connection_idle_timeout"`
	SensorUpdaterInterval    time.Duration `mapstructure:"sensor_updater_interval" yaml:"sensor_updater_interval"`
	StatusUpdaterInterval    time.Duration `mapstructure:"status_updater_interval" yaml:"status_updater_interval"`
	TopicDeviceType          mqtt.Topic    `mapstructure:"topic_device_type" yaml:"topic_device_type"`
	TopicRoomTemperature     mqtt.Topic    `mapstructure:"topic_room_temperature" yaml:"topic_room_temperature"`
	TopicFloorTemperature    mqtt.Topic    `mapstructure:"topic_floor_temperature" yaml:"topic_floor_temperature"`
	TopicHumidity            mqtt.Topic    `mapstructure:"topic_humidity" yaml:"topic_humidity"`
	TopicHeatingValve        mqtt.Topic    `mapstructure:"topic_heating_valve" yaml:"topic_heating_valve"`
	TopicCoolingValve        mqtt.Topic    `mapstructure:"topic_cooling_valve" yaml:"topic_cooling_valve"`
	TopicStatus              mqtt.Topic    `mapstructure:"topic_status" yaml:"topic_status"`
	TopicStatusState         mqtt.Topic    `mapstructure:"topic_status_state" yaml:"topic_status_state"`
	TopicHeatingOutputStatus mqtt.Topic    `mapstructure:"topic_heating_output_status" yaml:"topic_heating_output_status"`
	TopicHoldingFunction     mqtt.Topic    `mapstructure:"topic_holding_function" yaml:"topic_holding_function"`
	TopicFloorOverheat       mqtt.Topic    `mapstructure:"topic_floor_overheat" yaml:"topic_floor_overheat"`
	TopicFanSpeedNumbers     mqtt.Topic    `mapstructure:"topic_fan_speed_numbers" yaml:"topic_fan_speed_numbers"`
	//TopicTargetTemperature       mqtt.Topic `mapstructure:"topic_target_temperature" yaml:"topic_target_temperature"`
	//TopicTargetTemperatureState  mqtt.Topic `mapstructure:"topic_target_temperature_state" yaml:"topic_target_temperature_state"`
	//TopicAway                    mqtt.Topic `mapstructure:"topic_away" yaml:"topic_away"`
	//TopicAwayState               mqtt.Topic `mapstructure:"topic_away_state" yaml:"topic_away_state"`
	//TopicAwayTemperature         mqtt.Topic `mapstructure:"topic_away_temperature" yaml:"topic_away_temperature"`
	//TopicAwayTemperatureState    mqtt.Topic `mapstructure:"topic_away_temperature_state" yaml:"topic_away_temperature_state"`
	//TopicHoldingTemperature      mqtt.Topic `mapstructure:"topic_holding_temperature" yaml:"topic_holding_temperature"`
	//TopicHoldingTemperatureState mqtt.Topic `mapstructure:"topic_holding_temperature_state" yaml:"topic_holding_temperature_state"`
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
		StatusUpdaterInterval: time.Second * 30,

		TopicDeviceType:          prefix + "type",
		TopicRoomTemperature:     prefix + "temperature/room/state",
		TopicFloorTemperature:    prefix + "temperature/floor/state",
		TopicHumidity:            prefix + "humidity/state",
		TopicHeatingValve:        prefix + "valve/heating/state",
		TopicCoolingValve:        prefix + "valve/cooling/state",
		TopicStatus:              prefix + "power",
		TopicStatusState:         prefix + "power/state",
		TopicHeatingOutputStatus: prefix + "heating/state",
		TopicHoldingFunction:     prefix + "holding/state",
		TopicFloorOverheat:       prefix + "temperature/floor/overheat",
		TopicFanSpeedNumbers:     prefix + "fan/speed/numbers",
		//TopicTargetTemperature:       prefix + "target/temperature",
		//TopicTargetTemperatureState:  prefix + "target/temperature/state",
		//TopicAway:                    prefix + "away",
		//TopicAwayState:               prefix + "away/state",
		//TopicAwayTemperature:         prefix + "away/temperature",
		//TopicAwayTemperatureState:    prefix + "away/temperature/state",
		//TopicHoldingTemperature:      prefix + "holding/temperature",
		//TopicHoldingTemperatureState: prefix + "holding/temperature/state",
	}
}

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
	TopicHeatingOutputStatus mqtt.Topic    `mapstructure:"topic_heating_output_status" yaml:"topic_heating_output_status"`
	TopicRoomTemperature     mqtt.Topic    `mapstructure:"topic_room_temperature" yaml:"topic_room_temperature"`
	TopicFloorTemperature    mqtt.Topic    `mapstructure:"topic_floor_temperature" yaml:"topic_floor_temperature"`
	TopicHumidity            mqtt.Topic    `mapstructure:"topic_humidity" yaml:"topic_humidity"`
	TopicSetTemperature      mqtt.Topic    `mapstructure:"topic_set_temperature" yaml:"topic_set_temperature"`
	TopicSetTemperatureState mqtt.Topic    `mapstructure:"topic_set_temperature_state" yaml:"topic_set_temperature_state"`
}

func (t Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	var prefix mqtt.Topic = boggart.ComponentName + "/mc6/+/"

	return &Config{
		ProbesConfig:             probesConfig,
		LoggerConfig:             di.LoggerConfigDefaults(),
		ConnectionSlaveID:        0x1,
		ConnectionTimeout:        time.Second,
		ConnectionIdleTimeout:    time.Minute,
		SensorUpdaterInterval:    time.Minute,
		StatusUpdaterInterval:    time.Second * 10,
		TopicDeviceType:          prefix + "type",
		TopicHeatingOutputStatus: prefix + "status/output-heating",
		TopicRoomTemperature:     prefix + "room-temperature",
		TopicFloorTemperature:    prefix + "floor-temperature",
		TopicHumidity:            prefix + "humidity",
		TopicSetTemperature:      prefix + "temperature",
		TopicSetTemperatureState: prefix + "temperature/state",
	}
}

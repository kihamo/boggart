package zigbee2mqtt

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	NewAPI                  bool       `mapstructure:"new_api" yaml:"new_api"`
	TopicPrefix             mqtt.Topic `mapstructure:"topic_prefix" yaml:"topic_prefix"`
	TopicState              mqtt.Topic `mapstructure:"topic_state" yaml:"topic_state"`
	TopicConfig             mqtt.Topic `mapstructure:"topic_config" yaml:"topic_config"`
	TopicLog                mqtt.Topic `mapstructure:"topic_log" yaml:"topic_log"`
	TopicDevicesRequest     mqtt.Topic `mapstructure:"topic_devices_request" yaml:"topic_devices_request"`
	TopicDevicesResponse    mqtt.Topic `mapstructure:"topic_devices_response" yaml:"topic_devices_response"`
	TopicPermitJoin         mqtt.Topic `mapstructure:"topic_permit_join" yaml:"topic_permit_join"`
	TopicLastSeen           mqtt.Topic `mapstructure:"topic_last_seen" yaml:"topic_last_seen"`
	TopicElapsed            mqtt.Topic `mapstructure:"topic_elapsed" yaml:"topic_elapsed"`
	TopicReset              mqtt.Topic `mapstructure:"topic_reset" yaml:"topic_reset"`
	TopicFactoryReset       mqtt.Topic `mapstructure:"topic_factory_reset" yaml:"topic_factory_reset"`
	TopicLogLevel           mqtt.Topic `mapstructure:"topic_log_level" yaml:"topic_log_level"`
	TopicNetworkMapRequest  mqtt.Topic `mapstructure:"topic_networkmap_request" yaml:"topic_networkmap_request"`
	TopicNetworkMapResponse mqtt.Topic `mapstructure:"topic_networkmap_response" yaml:"topic_networkmap_response"`

	// New API
	TopicHealthCheckRequest  mqtt.Topic `mapstructure:"topic_health_check_request" yaml:"topic_health_check_request"`
	TopicHealthCheckResponse mqtt.Topic `mapstructure:"topic_health_check_response" yaml:"topic_health_check_response"`
	TopicLogging             mqtt.Topic `mapstructure:"topic_logging" yaml:"topic_logging"`
	TopicInfo                mqtt.Topic `mapstructure:"topic_info" yaml:"topic_info"`
	TopicDevices             mqtt.Topic `mapstructure:"topic_devices" yaml:"topic_devices"`
}

func (Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 30
	probesConfig.ReadinessTimeout = time.Second

	var prefix mqtt.Topic = "zigbee2mqtt"

	return &Config{
		ProbesConfig:            probesConfig,
		LoggerConfig:            di.LoggerConfigDefaults(),
		NewAPI:                  false,
		TopicPrefix:             prefix,
		TopicState:              prefix + "/bridge/state",
		TopicConfig:             prefix + "/bridge/config",
		TopicLog:                prefix + "/bridge/log",
		TopicDevicesRequest:     prefix + "/bridge/config/devices/get",
		TopicDevicesResponse:    prefix + "/bridge/config/devices",
		TopicPermitJoin:         prefix + "/bridge/config/permit_join",
		TopicLastSeen:           prefix + "/bridge/config/last_seen",
		TopicElapsed:            prefix + "/bridge/config/elapsed",
		TopicReset:              prefix + "/bridge/config/reset",
		TopicFactoryReset:       prefix + "/bridge/config/touchlink/factory_reset",
		TopicLogLevel:           prefix + "/bridge/config/log_level",
		TopicNetworkMapRequest:  prefix + "/bridge/networkmap",
		TopicNetworkMapResponse: prefix + "/bridge/networkmap/raw",

		TopicHealthCheckRequest:  prefix + "/bridge/request/health_check",
		TopicHealthCheckResponse: prefix + "/bridge/response/health_check",
		TopicLogging:             prefix + "/bridge/logging",
		TopicInfo:                prefix + "/bridge/info",
		TopicDevices:             prefix + "/bridge/devices",
	}
}

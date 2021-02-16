package esp

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	DeviceID                           string     `valid:"required" mapstructure:"device_id" yaml:"device_id"`
	BaseTopic                          string     `mapstructure:"base_topic" yaml:"base_topic"`
	TopicBroadcast                     mqtt.Topic `mapstructure:"topic_broadcast" yaml:"topic_broadcast"`
	TopicReset                         mqtt.Topic `mapstructure:"topic_reset" yaml:"topic_reset"`
	TopicRestart                       mqtt.Topic `mapstructure:"topic_restart" yaml:"topic_restart"`
	TopicDeviceAttribute               mqtt.Topic `mapstructure:"topic_device_attribute" yaml:"topic_device_attribute"`
	TopicDeviceAttributeFirmware       mqtt.Topic `mapstructure:"topic_device_attribute_firmware" yaml:"topic_device_attribute_firmware"`
	TopicDeviceAttributeImplementation mqtt.Topic `mapstructure:"topic_device_attribute_implementation" yaml:"topic_device_attribute_implementation"`
	TopicDeviceAttributeStats          mqtt.Topic `mapstructure:"topic_device_attribute_stats" yaml:"topic_device_attribute_stats"`
	TopicNodeAttribute                 mqtt.Topic `mapstructure:"topic_node_attribute" yaml:"topic_node_attribute"`
	TopicNodeList                      mqtt.Topic `mapstructure:"topic_node_list" yaml:"topic_node_list"`
	TopicNodeProperty                  mqtt.Topic `mapstructure:"topic_node_property" yaml:"topic_node_property"`
	TopicSettings                      mqtt.Topic `mapstructure:"topic_settings" yaml:"topic_settings"`
	TopicSettingsSet                   mqtt.Topic `mapstructure:"topic_settings_set" yaml:"topic_settings_set"`
	TopicOTAFirmware                   mqtt.Topic `mapstructure:"topic_ota_firmware" yaml:"topic_ota_firmware"`
	TopicOTAStatus                     mqtt.Topic `mapstructure:"topic_ota_status" yaml:"topic_ota_status"`
	TopicOTAEnabled                    mqtt.Topic `mapstructure:"topic_ota_enabled" yaml:"topic_ota_enabled"`
}

func (t Type) ConfigDefaults() interface{} {
	var (
		prefix     mqtt.Topic = "+/+/"
		prefixImpl            = prefix + "$implementation/"
	)

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 15
	probesConfig.ReadinessTimeout = time.Second

	return &Config{
		ProbesConfig:                       probesConfig,
		LoggerConfig:                       di.LoggerConfigDefaults(),
		BaseTopic:                          "homie",
		TopicBroadcast:                     "+/$broadcast/+",
		TopicReset:                         prefixImpl + "reset",
		TopicRestart:                       prefixImpl + "restart",
		TopicDeviceAttribute:               prefix + "+",
		TopicDeviceAttributeFirmware:       prefix + "$fw/+",
		TopicDeviceAttributeImplementation: prefixImpl + "+",
		TopicDeviceAttributeStats:          prefix + "$stats/+",
		TopicNodeAttribute:                 prefix + "$nodes",
		TopicNodeList:                      prefix + "+/+",
		TopicNodeProperty:                  prefix + "+/+/+",
		TopicSettings:                      prefixImpl + "config",
		TopicSettingsSet:                   prefixImpl + "config/set",
		TopicOTAFirmware:                   prefixImpl + "ota/firmware/+",
		TopicOTAStatus:                     prefixImpl + "ota/status",
		TopicOTAEnabled:                    prefixImpl + "ota/enabled",
	}
}

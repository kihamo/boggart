package root

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	DeviceIDFile       string     `mapstructure:"device_id_file" yaml:"device_id_file"`
	RuntimeConfigFile  string     `mapstructure:"runtime_config_file" yaml:"runtime_config_file"`
	TopicRuntimeConfig mqtt.Topic `mapstructure:"topic_runtime_config" yaml:"topic_runtime_config"`
}

func (t Type) Config() interface{} {
	return &Config{
		DeviceIDFile:       "/mnt/data/miio/device.uid",
		RuntimeConfigFile:  "/mnt/data/rockrobo/RoboController.cfg",
		TopicRuntimeConfig: boggart.ComponentName + "/+/runtime/+",
	}
}

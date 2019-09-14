package samsung_tizen

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Host                 string        `valid:"host,required"`
	LivenessInterval     time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout      time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	TopicPower           mqtt.Topic    `mapstructure:"topic_power" yaml:"topic_power"`
	TopicKey             mqtt.Topic    `mapstructure:"topic_key" yaml:"topic_key"`
	TopicDeviceID        mqtt.Topic    `mapstructure:"topic_device_id" yaml:"topic_device_id"`
	TopicDeviceModelName mqtt.Topic    `mapstructure:"topic_device_mode_name" yaml:"topic_device_mode_name"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/tv/+/"

	return &Config{
		LivenessInterval:     time.Second * 30,
		LivenessTimeout:      time.Second * 5,
		TopicPower:           prefix + "power",
		TopicKey:             prefix + "key",
		TopicDeviceID:        prefix + "device/id",
		TopicDeviceModelName: prefix + "device/model-name",
	}
}

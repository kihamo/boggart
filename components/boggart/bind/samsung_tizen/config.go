package samsung_tizen

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`

	Host                 string     `valid:"host,required"`
	TopicPower           mqtt.Topic `mapstructure:"topic_power" yaml:"topic_power"`
	TopicKey             mqtt.Topic `mapstructure:"topic_key" yaml:"topic_key"`
	TopicDeviceID        mqtt.Topic `mapstructure:"topic_device_id" yaml:"topic_device_id"`
	TopicDeviceModelName mqtt.Topic `mapstructure:"topic_device_mode_name" yaml:"topic_device_mode_name"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/tv/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 5,
		},
		TopicPower:           prefix + "power",
		TopicKey:             prefix + "key",
		TopicDeviceID:        prefix + "device/id",
		TopicDeviceModelName: prefix + "device/model-name",
	}
}

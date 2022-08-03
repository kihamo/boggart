package mqtt

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

	AES256Key              string `mapstructure:"aes256_key" yaml:"aes256_key" valid:"length(32|32),required"`
	MAC                    types.HardwareAddr
	TopicMC6Update         mqtt.Topic `mapstructure:"topic_mc6_update" yaml:"topic_mc6_update"`
	TopicMC6SetTemperature mqtt.Topic `mapstructure:"topic_mc6_set_temperature" yaml:"topic_mc6_set_temperature"`
	TopicSetTemperature    mqtt.Topic `mapstructure:"topic_set_temperature" yaml:"topic_set_temperature"`
}

func (t Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	var prefix mqtt.Topic = boggart.ComponentName + "/mc6/+/"

	return &Config{
		ProbesConfig:           probesConfig,
		LoggerConfig:           di.LoggerConfigDefaults(),
		TopicMC6Update:         "updData/+",
		TopicMC6SetTemperature: "+",
		TopicSetTemperature:    prefix + "temperature",
	}
}

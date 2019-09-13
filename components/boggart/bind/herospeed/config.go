package herospeed

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address                   boggart.URL   `valid:",required"`
	LivenessInterval          time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout           time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	PreviewRefreshInterval    time.Duration `mapstructure:"preview_refresh_interval" yaml:"preview_refresh_interval,omitempty"`
	TopicStateModel           mqtt.Topic    `mapstructure:"topic_state_model" yaml:"topic_state_model"`
	TopicStateFirmwareVersion mqtt.Topic    `mapstructure:"topic_state_firmware_version" yaml:"topic_state_firmware_version"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/state/"

	return &Config{
		LivenessInterval:          time.Minute,
		LivenessTimeout:           time.Second * 5,
		PreviewRefreshInterval:    time.Second * 5,
		TopicStateModel:           prefix + "model",
		TopicStateFirmwareVersion: prefix + "firmware/version",
	}
}

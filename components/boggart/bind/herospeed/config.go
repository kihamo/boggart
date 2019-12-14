package herospeed

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Address                   boggart.URL   `valid:",required"`
	PreviewRefreshInterval    time.Duration `mapstructure:"preview_refresh_interval" yaml:"preview_refresh_interval,omitempty"`
	TopicStateModel           mqtt.Topic    `mapstructure:"topic_state_model" yaml:"topic_state_model"`
	TopicStateFirmwareVersion mqtt.Topic    `mapstructure:"topic_state_firmware_version" yaml:"topic_state_firmware_version"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/state/"

	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		PreviewRefreshInterval:    time.Second * 5,
		TopicStateModel:           prefix + "model",
		TopicStateFirmwareVersion: prefix + "firmware/version",
	}
}

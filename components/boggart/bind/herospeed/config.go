package herospeed

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

	Address                   types.URL     `valid:",required"`
	PreviewRefreshInterval    time.Duration `mapstructure:"preview_refresh_interval" yaml:"preview_refresh_interval,omitempty"`
	TopicStateModel           mqtt.Topic    `mapstructure:"topic_state_model" yaml:"topic_state_model"`
	TopicStateFirmwareVersion mqtt.Topic    `mapstructure:"topic_state_firmware_version" yaml:"topic_state_firmware_version"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/state/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:              probesConfig,
		LoggerConfig:              di.LoggerConfigDefaults(),
		PreviewRefreshInterval:    time.Second * 5,
		TopicStateModel:           prefix + "model",
		TopicStateFirmwareVersion: prefix + "firmware/version",
	}
}

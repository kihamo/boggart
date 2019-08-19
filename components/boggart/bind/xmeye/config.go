package xmeye

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultLivenessInterval       = time.Minute
	DefaultLivenessTimeout        = time.Second * 5
	DefaultUpdaterInterval        = time.Minute
	DefaultUpdaterTimeout         = time.Second * 30
	DefaultAlarmStreamingInterval = time.Second * 5
)

type Config struct {
	Address                boggart.URL   `valid:",required"`
	LivenessInterval       time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout        time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval        time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout         time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	AlarmStreamingEnabled  bool          `mapstructure:"alarm_streaming_enabled" yaml:"alarm_streaming_enabled,omitempty"`
	AlarmStreamingInterval time.Duration `mapstructure:"alarm_streaming_interval" yaml:"alarm_streaming_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval:       DefaultLivenessInterval,
		LivenessTimeout:        DefaultLivenessTimeout,
		UpdaterInterval:        DefaultUpdaterInterval,
		UpdaterTimeout:         DefaultUpdaterTimeout,
		AlarmStreamingInterval: DefaultAlarmStreamingInterval,
	}
}

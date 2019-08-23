package hilink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultDebug            = false
	DefaultLivenessInterval = time.Minute
	DefaultLivenessTimeout  = time.Second * 5
	DefaultUpdaterInterval  = time.Minute
	DefaultUpdaterTimeout   = time.Second * 30
)

type Config struct {
	Address          boggart.URL `valid:",required"`
	Debug            bool
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout   time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		Debug:            DefaultDebug,
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
		UpdaterInterval:  DefaultUpdaterInterval,
		UpdaterTimeout:   DefaultUpdaterTimeout,
	}
}

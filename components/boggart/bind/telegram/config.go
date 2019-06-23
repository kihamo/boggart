package telegram

import (
	"time"
)

const (
	DefaultDebug            = false
	DefaultUpdatesEnabled   = false
	DefaultUpdatesBuffer    = 100
	DefaultUpdatesTimeout   = 60
	DefaultLivenessInterval = time.Minute
	DefaultLivenessTimeout  = time.Second * 5
)

type Config struct {
	Token            string
	Debug            bool
	UpdatesEnabled   bool          `mapstructure:"updates_enabled" yaml:"updates_enabled"`
	UpdatesBuffer    int           `mapstructure:"updates_buffer" yaml:"updates_buffer"`
	UpdatesTimeout   int           `mapstructure:"updates_timeout" yaml:"updates_timeout"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		Debug:            DefaultDebug,
		UpdatesEnabled:   DefaultUpdatesEnabled,
		UpdatesBuffer:    DefaultUpdatesBuffer,
		UpdatesTimeout:   DefaultUpdatesTimeout,
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
	}
}

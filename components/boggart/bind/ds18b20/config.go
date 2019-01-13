package ds18b20

import (
	"time"
)

const (
	DefaultLivenessInterval = time.Minute
	DefaultLivenessTimeout  = time.Second * 5
	DefaultUpdaterInterval  = time.Minute
)

type Config struct {
	Address          string        `valid:"required"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
		UpdaterInterval:  DefaultUpdaterInterval,
	}
}

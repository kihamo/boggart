package google_home

import (
	"time"
)

const (
	DefaultLivenessInterval = time.Second * 30
	DefaultLivenessTimeout  = time.Second * 10
)

type Config struct {
	Host             string        `valid:"host,required"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
	}
}

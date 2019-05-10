package miio

import (
	"time"
)

const (
	DefaultLivenessInterval = time.Second * 30
	DefaultLivenessTimeout  = time.Second * 5
)

type Config struct {
	Host             string        `valid:"host,required"`
	Token            string        `valid:"required"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
	}
}

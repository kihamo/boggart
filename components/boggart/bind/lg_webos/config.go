package lg_webos

import (
	"time"
)

const (
	DefaultLivenessInterval = time.Second * 30
	DefaultLivenessTimeout  = time.Second * 10
)

type Config struct {
	Host             string `valid:"host,required"`
	Key              string `valid:"required"`
	LivenessInterval string `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  string `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: DefaultLivenessInterval.String(),
		LivenessTimeout:  DefaultLivenessTimeout.String(),
	}
}

package esphome

import (
	"time"
)

type Config struct {
	Address          string `valid:"required"`
	Password         string
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: time.Minute,
		LivenessTimeout:  time.Second * 5,
	}
}

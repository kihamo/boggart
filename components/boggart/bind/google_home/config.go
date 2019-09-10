package google_home

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/google/home"
)

const (
	DefaultLivenessInterval = time.Second * 30
	DefaultLivenessTimeout  = time.Second * 10
)

type Config struct {
	Host             boggart.IP    `valid:",required"`
	Port             int           `valid:"port"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		Port:             home.DefaultPort,
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
	}
}

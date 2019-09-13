package google_home

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/google/home"
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
		LivenessInterval: time.Second * 30,
		LivenessTimeout:  time.Second * 10,
	}
}

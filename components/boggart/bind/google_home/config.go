package google_home

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/google/home"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Host boggart.IP `valid:",required"`
	Port int        `valid:"port"`
}

func (t Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 10,
		},
		Port: home.DefaultPort,
	}
}

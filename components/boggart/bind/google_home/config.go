package google_home

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Address boggart.URL `valid:",required"`
	Debug   bool
}

func (t Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 10,
		},
	}
}

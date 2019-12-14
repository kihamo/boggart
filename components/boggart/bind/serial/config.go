package serial

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Network string
	Host    string
	Port    int64
	Target  string
	Timeout time.Duration
}

func (t Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		Network: "tcp",
		Host:    "0.0.0.0",
		Port:    8600,
		Target:  "/dev/ttyUSB0",
		Timeout: time.Second,
	}
}

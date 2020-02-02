package openhab

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address boggart.URL `valid:",required"`
	Debug   bool
}

func (t Type) Config() interface{} {
	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
	}
}

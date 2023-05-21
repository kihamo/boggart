package keenetic

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address         types.URL `valid:",required"`
	Debug           bool
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) ConfigDefaults() interface{} {

	return &Config{
		ProbesConfig:    di.ProbesConfigDefaults(),
		LoggerConfig:    di.LoggerConfigDefaults(),
		UpdaterInterval: time.Minute * 5,
	}
}

package pass24online

import (
	"github.com/kihamo/boggart/components/boggart/di"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Phone    string `valid:",required"`
	Password string `valid:",required"`
	Debug    bool
}

func (t Type) Config() interface{} {
	return &Config{
		LoggerConfig: di.LoggerConfigDefaults(),
	}
}

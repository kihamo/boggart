package elektroset

import (
	"time"
)

const (
	DefaultUpdaterInterval = time.Hour
)

type Config struct {
	Login           string        `valid:"required"`
	Password        string        `valid:"required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	Debug           bool
}

func (Type) Config() interface{} {
	return &Config{
		UpdaterInterval: DefaultUpdaterInterval,
	}
}

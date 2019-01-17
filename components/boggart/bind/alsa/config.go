package alsa

import (
	"time"
)

const (
	DefaultUpdateInterval = time.Second * 3
)

type Config struct {
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		UpdaterInterval: DefaultUpdateInterval,
	}
}

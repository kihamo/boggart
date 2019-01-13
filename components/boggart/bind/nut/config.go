package nut

import (
	"time"
)

const (
	DefaultUpdaterInterval = time.Minute
)

type Config struct {
	Host            string `valid:"host,required"`
	Username        string
	Password        string
	UPS             string `valid:"required"`
	UpdaterInterval string `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		UpdaterInterval: DefaultUpdaterInterval.String(),
	}
}

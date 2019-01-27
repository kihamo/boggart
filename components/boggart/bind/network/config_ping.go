package network

import (
	"time"
)

const (
	DefaultRetry           = 1
	DefaultTimeout         = time.Second * 5
	DefaultUpdaterInterval = time.Minute
)

type ConfigPing struct {
	Hostname        string `valid:"host,required"`
	Retry           int
	Timeout         time.Duration
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t TypePing) Config() interface{} {
	return &ConfigPing{
		Retry:           DefaultRetry,
		Timeout:         DefaultTimeout,
		UpdaterInterval: DefaultUpdaterInterval,
	}
}

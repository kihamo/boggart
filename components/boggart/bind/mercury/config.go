package mercury

import (
	"time"
)

const (
	DefaultRS485Timeout    = time.Second
	DefaultUpdaterInterval = time.Minute
)

type Config struct {
	RS485Address    string        `mapstructure:"rs485_address" yaml:"rs485_address" valid:"required"`
	RS485Timeout    time.Duration `mapstructure:"rs485_timeout" yaml:"rs485_timeout"`
	Address         string        `valid:"required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		RS485Timeout:    DefaultRS485Timeout,
		UpdaterInterval: DefaultUpdaterInterval,
	}
}

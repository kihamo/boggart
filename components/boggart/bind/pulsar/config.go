package pulsar

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
	Address         string
	Input1Offset    float32       `mapstructure:"input1_offset" valid:"float"`
	Input2Offset    float32       `mapstructure:"input2_offset" valid:"float"`
	Input3Offset    float32       `mapstructure:"input3_offset" valid:"float"`
	Input4Offset    float32       `mapstructure:"input4_offset" valid:"float"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		RS485Timeout:    DefaultRS485Timeout,
		UpdaterInterval: DefaultUpdaterInterval,
	}
}

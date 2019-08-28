package v1

import (
	"time"
)

const (
	DefaultRS485Timeout = time.Second

	/*
		При отсутствии тока в последовательной цепи и значении напряжения, равном 1,15Uном, испытательный выход
		счётчика не создаёт более одного импульса в течение времени, равного 4,4 мин и 3,5 мин для счётчиков класса
		точности 1 и 2 соответственно.
	*/
	DefaultUpdaterInterval = time.Minute * 5
)

type Config struct {
	RS485Address    string        `mapstructure:"rs485_address" yaml:"rs485_address" valid:"required"`
	RS485Timeout    time.Duration `mapstructure:"rs485_timeout" yaml:"rs485_timeout"`
	Address         string        `valid:"required"`
	Location        string
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		RS485Timeout:    DefaultRS485Timeout,
		Location:        time.Now().Location().String(),
		UpdaterInterval: DefaultUpdaterInterval,
	}
}

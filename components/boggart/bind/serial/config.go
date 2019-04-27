package serial

import (
	"time"
)

const (
	DefaultTimeout = time.Second
	DefaultHost    = "0.0.0.0"
	DefaultPort    = 8600
)

type Config struct {
	Host    string `valid:"required"`
	Port    int64  `valid:"required"`
	Target  string `valid:"required"`
	Timeout time.Duration
}

func (t Type) Config() interface{} {
	return &Config{
		Host:    DefaultHost,
		Port:    DefaultPort,
		Timeout: DefaultTimeout,
	}
}

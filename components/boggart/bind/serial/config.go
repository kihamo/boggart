package serial

import (
	"time"
)

type Config struct {
	Network string
	Host    string
	Port    int64
	Target  string
	Timeout time.Duration
}

func (t Type) Config() interface{} {
	return &Config{
		Network: "tcp",
		Host:    "0.0.0.0",
		Port:    8600,
		Target:  "/dev/ttyUSB0",
		Timeout: time.Second,
	}
}

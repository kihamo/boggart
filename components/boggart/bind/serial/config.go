package serial

import (
	"time"

	"github.com/kihamo/boggart/protocols/serial"
)

const (
	DefaultNetwork = "tcp"
	DefaultHost    = "0.0.0.0"
	DefaultPort    = 8600
	DefaultTarget  = serial.DefaultSerialAddress
	DefaultTimeout = time.Second
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
		Network: DefaultNetwork,
		Host:    DefaultHost,
		Port:    DefaultPort,
		Target:  DefaultTarget,
		Timeout: DefaultTimeout,
	}
}

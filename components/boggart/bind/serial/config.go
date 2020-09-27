package serial

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/serial"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Network  string
	Host     string
	Port     int64
	Target   string
	BaudRate int `mapstructure:"baud_rate" yaml:"baud_rate"`
	DataBits int `mapstructure:"data_bits" yaml:"data_bits"`
	StopBits int `mapstructure:"stop_bits" yaml:"stop_bits"`
	Parity   string
	Timeout  time.Duration
	Once     bool
}

func (t Type) Config() interface{} {
	serialConfig := serial.DefaultOptions()

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		Network:      "tcp",
		Host:         "0.0.0.0",
		Port:         8600,
		Target:       serialConfig.Address,
		BaudRate:     serialConfig.BaudRate,
		DataBits:     serialConfig.DataBits,
		StopBits:     serialConfig.StopBits,
		Parity:       serialConfig.Parity,
		Timeout:      serialConfig.Timeout,
		Once:         true,
	}
}

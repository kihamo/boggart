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
	def := serial.DefaultOptions()

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfigDefaults(),
		Network:      "tcp",
		Host:         "0.0.0.0",
		Port:         8600,
		Target:       def.Address,
		BaudRate:     def.BaudRate,
		DataBits:     def.DataBits,
		StopBits:     def.StopBits,
		Parity:       def.Parity,
		Timeout:      def.Timeout,
		Once:         true,
	}
}

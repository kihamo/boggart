package openhab

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address      types.URL `valid:",required"`
	Debug        bool
	ProxyEnabled bool   `mapstructure:"proxy_enabled" yaml:"proxy_enabled"`
	ProxyAddress string `mapstructure:"proxy_address" yaml:"proxy_address"`
}

func (t Type) Config() interface{} {
	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		ProxyEnabled: true,
		ProxyAddress: ":8089",
	}
}

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

func (t Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 30
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig: probesConfig,
		LoggerConfig: di.LoggerConfigDefaults(),
		ProxyEnabled: true,
		ProxyAddress: ":8089",
	}
}

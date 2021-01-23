package dom24

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Accounts                []string
	Phone                   string `valid:"required"`
	Password                string `valid:"required"`
	AutoRegisterIfNotExists bool   `mapstructure:"auto_register_if_not_exists" yaml:"auto_register_if_not_exists"`
	Debug                   bool
}

func (Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Hour
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig:            probesConfig,
		LoggerConfig:            di.LoggerConfigDefaults(),
		AutoRegisterIfNotExists: true,
		Debug:                   false,
	}
}

package cloud

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	ProviderLink    types.URL `mapstructure:"provider_link" yaml:"provider_link"  valid:"required"`
	Login           string    `valid:",required"`
	ApiKey          string    `mapstructure:"api_key" yaml:"api_key" valid:",required"`
	DeviceID        int64     `mapstructure:"device_id" yaml:"device_id"`
	Debug           bool
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout  time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	/*
		DS датчики температуры опрашиваются раз в 10 минут
		Остальные кажется что раз в минуту, но это не точно
	*/
}

func (t Type) ConfigDefaults() interface{} {
	cfg := &Config{
		ProbesConfig:    di.ProbesConfigDefaults(),
		LoggerConfig:    di.LoggerConfigDefaults(),
		UpdaterInterval: time.Minute,
	}

	if t.Link != nil {
		cfg.ProviderLink.URL = *t.Link
	}

	return cfg
}

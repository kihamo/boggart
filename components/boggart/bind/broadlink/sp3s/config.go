package sp3s

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultUpdateInterval    = time.Second * 3 // as e-control app, refresh every 3 sec
	DefaultConnectionTimeout = time.Second
)

type Config struct {
	Host              string               `valid:",required"`
	MAC               boggart.HardwareAddr `valid:",required"`
	Model             string               `valid:"in(sp3seu|sp3sus),required"`
	UpdaterInterval   time.Duration        `mapstructure:"updater_interval" yaml:"updater_interval"`
	ConnectionTimeout time.Duration        `mapstructure:"connection_timeout" yaml:"connection_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		UpdaterInterval:   DefaultUpdateInterval,
		ConnectionTimeout: DefaultConnectionTimeout,
	}
}

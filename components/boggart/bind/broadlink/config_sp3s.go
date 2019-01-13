package broadlink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	SP3SDefaultUpdateInterval = time.Second * 3 // as e-control app, refresh every 3 sec
)

type ConfigSP3S struct {
	IP              boggart.IP           `valid:",required"`
	MAC             boggart.HardwareAddr `valid:",required"`
	Model           string               `valid:"in(sp3seu|sp3sus),required"`
	UpdaterInterval time.Duration        `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t TypeSP3S) Config() interface{} {
	return &ConfigSP3S{
		UpdaterInterval: SP3SDefaultUpdateInterval,
	}
}

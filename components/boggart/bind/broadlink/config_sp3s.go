package broadlink

import (
	"time"
)

const (
	SP3SDefaultUpdateInterval = time.Second * 3 // as e-control app, refresh every 3 sec
)

type ConfigSP3S struct {
	IP              string        `valid:"ip,required"`
	MAC             string        `valid:"mac,required"`
	Model           string        `valid:"in(sp3seu|sp3sus),required"`
	UpdaterInterval time.Duration `mapstructure:"updater_interval"`
}

func (t TypeSP3S) Config() interface{} {
	return &ConfigSP3S{
		UpdaterInterval: SP3SDefaultUpdateInterval,
	}
}

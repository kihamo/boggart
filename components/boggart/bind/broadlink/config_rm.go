package broadlink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	RMCaptureDuration = time.Second * 15
)

type ConfigRM struct {
	IP              boggart.IP           `valid:",required"`
	MAC             boggart.HardwareAddr `valid:",required"`
	Model           string               `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration time.Duration        `mapstructure:"capture_interval" yaml:"capture_interval"`
}

func (t TypeRM) Config() interface{} {
	return &ConfigRM{
		CaptureDuration: RMCaptureDuration,
	}
}

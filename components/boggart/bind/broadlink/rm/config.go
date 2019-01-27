package rm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultCaptureDuration  = time.Second * 15
	DefaultLivenessInterval = time.Second * 30
	DefaultLivenessTimeout  = time.Second * 10
)

type ConfigRM struct {
	IP               boggart.IP           `valid:",required"`
	MAC              boggart.HardwareAddr `valid:",required"`
	Model            string               `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration  time.Duration        `mapstructure:"capture_interval" yaml:"capture_interval"`
	LivenessInterval time.Duration        `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration        `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &ConfigRM{
		CaptureDuration:  DefaultCaptureDuration,
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
	}
}

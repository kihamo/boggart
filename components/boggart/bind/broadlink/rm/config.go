package rm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultCaptureDuration   = time.Second * 15
	DefaultLivenessInterval  = time.Second * 30
	DefaultLivenessTimeout   = time.Second * 10
	DefaultConnectionTimeout = time.Second
)

type ConfigRM struct {
	Host              string               `valid:",required"`
	MAC               boggart.HardwareAddr `valid:",required"`
	Model             string               `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration   time.Duration        `mapstructure:"capture_interval" yaml:"capture_interval"`
	LivenessInterval  time.Duration        `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout   time.Duration        `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	ConnectionTimeout time.Duration        `mapstructure:"connection_timeout" yaml:"connection_timeout"`
}

func (t Type) Config() interface{} {
	return &ConfigRM{
		CaptureDuration:   DefaultCaptureDuration,
		LivenessInterval:  DefaultLivenessInterval,
		LivenessTimeout:   DefaultLivenessTimeout,
		ConnectionTimeout: DefaultConnectionTimeout,
	}
}

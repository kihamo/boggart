package boggart

import (
	"time"
)

const (
	ReadinessProbeDefaultPeriod = time.Minute
	LivenessProbeDefaultPeriod  = time.Minute
)

type BindConfigReadinessProbePeriod interface {
	ReadinessProbePeriod() time.Duration
}

type BindConfigLivenessProbePeriod interface {
	LivenessProbePeriod() time.Duration
}

type BindConfig struct {
	ReadinessPeriod time.Duration `mapstructure:"readiness_probe_period" yaml:"readiness_probe_period"`
	LivenessPeriod  time.Duration `mapstructure:"liveness_probe_period" yaml:"liveness_probe_period"`
}

func (c BindConfig) ReadinessProbePeriod() time.Duration {
	if c.ReadinessPeriod <= 0 {
		return ReadinessProbeDefaultPeriod
	}

	return c.ReadinessPeriod
}

func (c BindConfig) LivenessProbePeriod() time.Duration {
	if c.LivenessPeriod <= 0 {
		return LivenessProbeDefaultPeriod
	}

	return c.LivenessPeriod
}

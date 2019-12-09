package boggart

import (
	"time"
)

const (
	ReadinessProbeDefaultPeriod  = time.Minute
	ReadinessProbeDefaultTimeout = time.Second * 30
	LivenessProbeDefaultPeriod   = time.Minute
	LivenessProbeDefaultTimeout  = time.Second * 30
)

type BindConfigReadinessProbe interface {
	ReadinessProbePeriod() time.Duration
	ReadinessProbeTimeout() time.Duration
}

type BindConfigLivenessProbe interface {
	LivenessProbePeriod() time.Duration
	LivenessProbeTimeout() time.Duration
}

type BindConfig struct {
	ReadinessPeriod  time.Duration `mapstructure:"readiness_probe_period" yaml:"readiness_probe_period"`
	ReadinessTimeout time.Duration `mapstructure:"readiness_probe_timeout" yaml:"readiness_probe_timeout"`
	LivenessPeriod   time.Duration `mapstructure:"liveness_probe_period" yaml:"liveness_probe_period"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_probe_timeout" yaml:"liveness_probe_timeout"`
}

func (c BindConfig) ReadinessProbePeriod() time.Duration {
	if c.ReadinessPeriod <= 0 {
		return ReadinessProbeDefaultPeriod
	}

	return c.ReadinessPeriod
}

func (c BindConfig) ReadinessProbeTimeout() time.Duration {
	return c.ReadinessTimeout
}

func (c BindConfig) LivenessProbePeriod() time.Duration {
	if c.LivenessPeriod <= 0 {
		return LivenessProbeDefaultPeriod
	}

	return c.LivenessPeriod
}

func (c BindConfig) LivenessProbeTimeout() time.Duration {
	return c.LivenessTimeout
}

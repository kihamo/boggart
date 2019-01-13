package mikrotik

import (
	"time"
)

const (
	DefaultLivenessInterval = time.Minute
	DefaultLivenessTimeout  = time.Second * 5
	DefaultUpdaterInterval  = time.Minute * 5
)

type Config struct {
	Address          string        `valid:"url,required"`
	SyslogClient     string        `valid:"url" mapstructure:"syslog_client" yaml:"syslog_client,omitempty"`
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
		UpdaterInterval:  DefaultUpdaterInterval,
	}
}

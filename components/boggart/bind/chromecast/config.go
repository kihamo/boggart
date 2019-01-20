package chromecast

import (
	"time"

	"github.com/barnybug/go-cast/log"
	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultPort             = 8009
	DefaultName             = boggart.ComponentName
	DefaultLivenessInterval = time.Second * 30
	DefaultLivenessTimeout  = time.Second * 10
)

type Config struct {
	Debug            bool
	Host             boggart.IP `valid:",required"`
	Port             int        `valid:"port"`
	Name             string
	LivenessInterval time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		Debug:            log.Debug,
		Port:             DefaultPort,
		Name:             DefaultName,
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
	}
}

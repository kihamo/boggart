package herospeed

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultLivenessInterval       = time.Minute
	DefaultLivenessTimeout        = time.Second * 5
	DefaultPreviewRefreshInterval = time.Second * 5
)

type Config struct {
	Address                boggart.URL   `valid:",required"`
	LivenessInterval       time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout        time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	PreviewRefreshInterval time.Duration `mapstructure:"preview_refresh_interval" yaml:"preview_refresh_interval,omitempty"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval:       DefaultLivenessInterval,
		LivenessTimeout:        DefaultLivenessTimeout,
		PreviewRefreshInterval: DefaultPreviewRefreshInterval,
	}
}

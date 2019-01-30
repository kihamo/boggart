package hikvision

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultLivenessInterval = time.Minute
	DefaultLivenessTimeout  = time.Second * 5
	DefaultUpdaterInterval  = time.Minute
	DefaultUpdaterTimeout   = time.Second * 30
	DefaultPTZInterval      = time.Minute
	DefaultPTZTimeout       = time.Second * 5
	DefaultWidgetChannel    = 101
)

type Config struct {
	Address              boggart.URL   `valid:",required"`
	LivenessInterval     time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout      time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval      time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout       time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	PTZInterval          time.Duration `mapstructure:"ptz_interval" yaml:"ptz_interval"`
	PTZTimeout           time.Duration `mapstructure:"ptz_timeout" yaml:"ptz_timeout"`
	PTZEnabled           bool          `mapstructure:"ptz_enabled" yaml:"ptz_enabled,omitempty"`
	EventsEnabled        bool          `mapstructure:"events_enabled" yaml:"events_enabled,omitempty"`
	EventsIgnoreInterval time.Duration `mapstructure:"events_ignore_interval" yaml:"events_ignore_interval,omitempty"`
	WidgetChannel        uint64        `mapstructure:"widget_channel" yaml:"widget_channel,omitempty"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval: DefaultLivenessInterval,
		LivenessTimeout:  DefaultLivenessTimeout,
		UpdaterInterval:  DefaultUpdaterInterval,
		UpdaterTimeout:   DefaultUpdaterTimeout,
		PTZInterval:      DefaultPTZInterval,
		PTZTimeout:       DefaultPTZTimeout,
		WidgetChannel:    DefaultWidgetChannel,
	}
}

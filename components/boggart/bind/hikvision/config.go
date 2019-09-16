package hikvision

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address                        boggart.URL   `valid:",required"`
	LivenessInterval               time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout                time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval                time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout                 time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	PTZInterval                    time.Duration `mapstructure:"ptz_interval" yaml:"ptz_interval"`
	PTZTimeout                     time.Duration `mapstructure:"ptz_timeout" yaml:"ptz_timeout"`
	PTZEnabled                     bool          `mapstructure:"ptz_enabled" yaml:"ptz_enabled,omitempty"`
	EventsEnabled                  bool          `mapstructure:"events_enabled" yaml:"events_enabled,omitempty"`
	EventsStreamingEnabled         bool          `mapstructure:"events_streaming_enabled" yaml:"events_streaming_enabled,omitempty"`
	EventsIgnoreInterval           time.Duration `mapstructure:"events_ignore_interval" yaml:"events_ignore_interval,omitempty"`
	WidgetChannel                  uint64        `mapstructure:"widget_channel" yaml:"widget_channel,omitempty"`
	PreviewRefreshInterval         time.Duration `mapstructure:"preview_refresh_interval" yaml:"preview_refresh_interval,omitempty"`
	TopicPTZMove                   mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZAbsolute               mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZContinuous             mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZRelative               mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZPreset                 mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZMomentary              mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZStatusElevation        mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZStatusAzimuth          mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicPTZStatusZoom             mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicEvent                     mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateModel                mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateFirmwareVersion      mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateFirmwareReleasedDate mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateUpTime               mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateMemoryUsage          mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateMemoryAvailable      mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateHDDCapacity          mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateHDDFree              mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
	TopicStateHDDUsage             mqtt.Topic    `mapstructure:"topic_" yaml:"topic_"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/"

	return &Config{
		LivenessInterval:               time.Minute,
		LivenessTimeout:                time.Second * 5,
		UpdaterInterval:                time.Minute,
		UpdaterTimeout:                 time.Second * 30,
		PTZInterval:                    time.Minute,
		PTZTimeout:                     time.Second * 5,
		EventsIgnoreInterval:           time.Second * 5,
		WidgetChannel:                  101,
		PreviewRefreshInterval:         time.Second * 5,
		TopicPTZMove:                   prefix + "ptz/+/move",
		TopicPTZAbsolute:               prefix + "ptz/+/absolute",
		TopicPTZContinuous:             prefix + "ptz/+/continuous",
		TopicPTZRelative:               prefix + "ptz/+/relative",
		TopicPTZPreset:                 prefix + "ptz/+/preset",
		TopicPTZMomentary:              prefix + "ptz/+/momentary",
		TopicPTZStatusElevation:        prefix + "ptz/+/status/elevation",
		TopicPTZStatusAzimuth:          prefix + "ptz/+/status/azimuth",
		TopicPTZStatusZoom:             prefix + "ptz/+/status/zoom",
		TopicEvent:                     prefix + "event/+/+",
		TopicStateModel:                prefix + "state/model",
		TopicStateFirmwareVersion:      prefix + "state/firmware/version",
		TopicStateFirmwareReleasedDate: prefix + "state/firmware/release-date",
		TopicStateUpTime:               prefix + "state/uptime",
		TopicStateMemoryUsage:          prefix + "state/memory/usage",
		TopicStateMemoryAvailable:      prefix + "state/memory/available",
		TopicStateHDDCapacity:          prefix + "state/hdd/+/capacity",
		TopicStateHDDFree:              prefix + "state/hdd/+/free",
		TopicStateHDDUsage:             prefix + "state/hdd/+/usage",
	}
}

package hikvision

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address                          types.URL     `valid:",required"`
	SystemTimeNTPAddress             types.URL     `mapstructure:"system_time_ntp_address" yaml:"system_time_ntp_address,omitempty"`
	TopicPTZAbsolute                 mqtt.Topic    `mapstructure:"topic_ptz_absolute" yaml:"topic_ptz_absolute"`
	TopicPTZContinuous               mqtt.Topic    `mapstructure:"topic_ptz_continuous" yaml:"topic_ptz_continuous"`
	TopicPTZRelative                 mqtt.Topic    `mapstructure:"topic_ptz_relative" yaml:"topic_ptz_relative"`
	TopicPTZPreset                   mqtt.Topic    `mapstructure:"topic_ptz_preset" yaml:"topic_ptz_preset"`
	TopicPTZMomentary                mqtt.Topic    `mapstructure:"topic_ptz_momentary" yaml:"topic_ptz_momentary"`
	TopicPTZStatusElevation          mqtt.Topic    `mapstructure:"topic_ptz_status_elevation" yaml:"topic_ptz_status_elevation"`
	TopicPTZStatusAzimuth            mqtt.Topic    `mapstructure:"topic_ptz_status_azimuth" yaml:"topic_ptz_status_azimuth"`
	TopicPTZStatusZoom               mqtt.Topic    `mapstructure:"topic_ptz_status_zoom" yaml:"topic_ptz_status_zoom"`
	TopicEvent                       mqtt.Topic    `mapstructure:"topic_event" yaml:"topic_event"`
	TopicStateModel                  mqtt.Topic    `mapstructure:"topic_state_model" yaml:"topic_state_model"`
	TopicStateFirmwareVersion        mqtt.Topic    `mapstructure:"topic_state_firmware_version" yaml:"topic_state_firmware_version"`
	TopicStateFirmwareReleasedDate   mqtt.Topic    `mapstructure:"topic_state_firmware_released_date" yaml:"topic_state_firmware_released_date"`
	TopicStateUpTime                 mqtt.Topic    `mapstructure:"topic_state_up_time" yaml:"topic_state_up_time"`
	TopicStateMemoryUsage            mqtt.Topic    `mapstructure:"topic_state_memory_usage" yaml:"topic_state_memory_usage"`
	TopicStateMemoryAvailable        mqtt.Topic    `mapstructure:"topic_state_available" yaml:"topic_state_available"`
	TopicStateHDDStatus              mqtt.Topic    `mapstructure:"topic_state_hdd_status" yaml:"topic_state_hdd_status"`
	TopicStateHDDCapacity            mqtt.Topic    `mapstructure:"topic_state_hdd_capacity" yaml:"topic_state_hdd_capacity"`
	TopicStateHDDFree                mqtt.Topic    `mapstructure:"topic_state_hdd_free" yaml:"topic_state_hdd_free"`
	TopicStateHDDUsage               mqtt.Topic    `mapstructure:"topic_state_hdd_usage" yaml:"topic_state_hdd_usage"`
	UpdaterInterval                  time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout                   time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	PTZInterval                      time.Duration `mapstructure:"ptz_interval" yaml:"ptz_interval"`
	PTZTimeout                       time.Duration `mapstructure:"ptz_timeout" yaml:"ptz_timeout"`
	EventsIgnoreInterval             time.Duration `mapstructure:"events_ignore_interval" yaml:"events_ignore_interval,omitempty"`
	PreviewRefreshInterval           time.Duration `mapstructure:"preview_refresh_interval" yaml:"preview_refresh_interval,omitempty"`
	WidgetChannel                    uint64        `mapstructure:"widget_channel" yaml:"widget_channel,omitempty"`
	SystemTimeNTPSynchronizeInterval uint64        `mapstructure:"system_time_ntp_synchronize_interval" yaml:"system_time_ntp_synchronize_interval,omitempty"`
	Debug                            bool
	EventsEnabled                    bool `mapstructure:"events_enabled" yaml:"events_enabled,omitempty"`
	EventsStreamingEnabled           bool `mapstructure:"events_streaming_enabled" yaml:"events_streaming_enabled,omitempty"`
	VirtualHostAutoEnabled           bool `mapstructure:"virtual_host_auto_enabled" yaml:"virtual_host_auto_enabled,omitempty"`
	SystemTimeNTPAutoEnabled         bool `mapstructure:"system_time_ntp_auto_enabled" yaml:"system_time_ntp_auto_enabled,omitempty"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/cctv/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		SystemTimeNTPAddress:             defaultNTPServer,
		SystemTimeNTPSynchronizeInterval: 60,
		ProbesConfig:                     probesConfig,
		LoggerConfig:                     di.LoggerConfigDefaults(),
		UpdaterInterval:                  time.Minute,
		UpdaterTimeout:                   time.Second * 30,
		PTZInterval:                      time.Minute,
		PTZTimeout:                       time.Second * 5,
		EventsIgnoreInterval:             time.Second * 5,
		WidgetChannel:                    101,
		PreviewRefreshInterval:           time.Second * 5,
		TopicPTZAbsolute:                 prefix + "ptz/+/absolute",
		TopicPTZContinuous:               prefix + "ptz/+/continuous",
		TopicPTZRelative:                 prefix + "ptz/+/relative",
		TopicPTZPreset:                   prefix + "ptz/+/preset",
		TopicPTZMomentary:                prefix + "ptz/+/momentary",
		TopicPTZStatusElevation:          prefix + "ptz/+/status/elevation",
		TopicPTZStatusAzimuth:            prefix + "ptz/+/status/azimuth",
		TopicPTZStatusZoom:               prefix + "ptz/+/status/zoom",
		TopicEvent:                       prefix + "event/+/+",
		TopicStateModel:                  prefix + "state/model",
		TopicStateFirmwareVersion:        prefix + "state/firmware/version",
		TopicStateFirmwareReleasedDate:   prefix + "state/firmware/release-date",
		TopicStateUpTime:                 prefix + "state/uptime",
		TopicStateMemoryUsage:            prefix + "state/memory/usage",
		TopicStateMemoryAvailable:        prefix + "state/memory/available",
		TopicStateHDDStatus:              prefix + "state/hdd/+/status",
		TopicStateHDDCapacity:            prefix + "state/hdd/+/capacity",
		TopicStateHDDFree:                prefix + "state/hdd/+/free",
		TopicStateHDDUsage:               prefix + "state/hdd/+/usage",
	}
}

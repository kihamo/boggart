package keenetic

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

	Address              types.URL `valid:",required"`
	Debug                bool
	HotspotSyncInterval  time.Duration `mapstructure:"hotspot_sync_interval" yaml:"hotspot_sync_interval"`
	UpdaterInterval      time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	SyncAfterSyslogDelay time.Duration `mapstructure:"sync_after_syslog_delay" yaml:"sync_after_syslog_delay"`
	TopicHotspotState    mqtt.Topic    `mapstructure:"topic_hotspot_state" yaml:"topic_hotspot_state"`
	TopicSyslog          mqtt.Topic    `mapstructure:"topic_syslog" yaml:"topic_syslog"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/keenetic/+/"

	return &Config{
		ProbesConfig:         di.ProbesConfigDefaults(),
		LoggerConfig:         di.LoggerConfigDefaults(),
		HotspotSyncInterval:  time.Minute,
		UpdaterInterval:      time.Minute * 5,
		SyncAfterSyslogDelay: time.Second,
		TopicHotspotState:    prefix + "hotspot/+/state",
	}
}

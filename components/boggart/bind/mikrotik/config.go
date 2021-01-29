package mikrotik

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Address                       string        `valid:"url,required"`
	ClientTimeout                 time.Duration `mapstructure:"client_timeout" yaml:"client_timeout"`
	ClientsSyncInterval           time.Duration `mapstructure:"clients_sync_interval" yaml:"clients_sync_interval"`
	UpdaterInterval               time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	SyslogTagWireless             string        `mapstructure:"topic_syslog_tag_wireless" yaml:"topic_syslog_tag_wireless"`
	SyslogTagL2TP                 string        `mapstructure:"topic_syslog_tag_l2tp" yaml:"topic_syslog_tag_l2tp"`
	TopicWiFiMACState             mqtt.Topic    `mapstructure:"topic_wifi_mac_state" yaml:"topic_wifi_mac_state"`
	TopicInterfaceConnect         mqtt.Topic    `mapstructure:"topic_interface_connect" yaml:"topic_interface_connect"`
	TopicPackagesInstalledVersion mqtt.Topic    `mapstructure:"topic_packages_installed_version" yaml:"topic_packages_installed_version"`
	TopicPackagesLatestVersion    mqtt.Topic    `mapstructure:"topic_packages_latest_version" yaml:"topic_packages_latest_version"`
	TopicFirmwareInstalledVersion mqtt.Topic    `mapstructure:"topic_firmware_installed_version" yaml:"topic_firmware_installed_version"`
	TopicFirmwareLatestVersion    mqtt.Topic    `mapstructure:"topic_firmware_latest_version" yaml:"topic_firmware_latest_version"`
	TopicSyslog                   mqtt.Topic    `mapstructure:"topic_syslog" yaml:"topic_syslog"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/router/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:                  probesConfig,
		LoggerConfig:                  di.LoggerConfigDefaults(),
		ClientTimeout:                 time.Second * 10,
		ClientsSyncInterval:           time.Minute,
		UpdaterInterval:               time.Minute * 5,
		SyslogTagWireless:             "wifi",
		SyslogTagL2TP:                 "vpn",
		TopicInterfaceConnect:         prefix + "+/+",
		TopicPackagesInstalledVersion: prefix + "packages/installed-version",
		TopicPackagesLatestVersion:    prefix + "packages/latest-version",
		TopicFirmwareInstalledVersion: prefix + "firmware/installed-version",
		TopicFirmwareLatestVersion:    prefix + "firmware/latest-version",
	}
}

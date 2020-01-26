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
	SyslogClient                  string        `valid:"url" mapstructure:"syslog_client" yaml:"syslog_client,omitempty"`
	ClientsSyncInterval           time.Duration `mapstructure:"clients_sync_interval" yaml:"clients_sync_interval"`
	UpdaterInterval               time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicWiFiMACState             mqtt.Topic    `mapstructure:"topic_wifi_mac_state" yaml:"topic_wifi_mac_state"`
	TopicWiFiConnectedMAC         mqtt.Topic    `mapstructure:"topic_wifi_connected_mac" yaml:"topic_wifi_connected_mac"`
	TopicWiFiDisconnectedMAC      mqtt.Topic    `mapstructure:"topic_wifi_disconnected_mac" yaml:"topic_wifi_disconnected_mac"`
	TopicVPNLoginState            mqtt.Topic    `mapstructure:"topic_vpn_login_state" yaml:"topic_vpn_login_state"`
	TopicVPNConnectedLogin        mqtt.Topic    `mapstructure:"topic_vpn_connected_login" yaml:"topic_vpn_connected_login"`
	TopicVPNDisconnectedLogin     mqtt.Topic    `mapstructure:"topic_vpn_disconnected_login" yaml:"topic_vpn_disconnected_login"`
	TopicPackagesInstalledVersion mqtt.Topic    `mapstructure:"topic_packages_installed_version" yaml:"topic_packages_installed_version"`
	TopicPackagesLatestVersion    mqtt.Topic    `mapstructure:"topic_packages_latest_version" yaml:"topic_packages_latest_version"`
	TopicFirmwareInstalledVersion mqtt.Topic    `mapstructure:"topic_firmware_installed_version" yaml:"topic_firmware_installed_version"`
	TopicFirmwareLatestVersion    mqtt.Topic    `mapstructure:"topic_firmware_latest_version" yaml:"topic_firmware_latest_version"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/router/+/"

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		ClientsSyncInterval:           time.Minute,
		UpdaterInterval:               time.Minute * 5,
		TopicWiFiMACState:             prefix + "wifi/clients/+/state",
		TopicWiFiConnectedMAC:         prefix + "wifi/clients/last/on/mac",
		TopicWiFiDisconnectedMAC:      prefix + "wifi/clients/last/off/mac",
		TopicVPNLoginState:            prefix + "vpn/clients/+/state",
		TopicVPNConnectedLogin:        prefix + "vpn/clients/last/on/login",
		TopicVPNDisconnectedLogin:     prefix + "vpn/clients/last/off/login",
		TopicPackagesInstalledVersion: prefix + "packages/installed-version",
		TopicPackagesLatestVersion:    prefix + "packages/latest-version",
		TopicFirmwareInstalledVersion: prefix + "firmware/installed-version",
		TopicFirmwareLatestVersion:    prefix + "firmware/latest-version",
	}
}

package mikrotik

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

	Address                       types.URL         `valid:",required"`
	ClientTimeout                 time.Duration     `mapstructure:"client_timeout" yaml:"client_timeout"`
	ClientsSyncInterval           time.Duration     `mapstructure:"clients_sync_interval" yaml:"clients_sync_interval"`
	UpdaterInterval               time.Duration     `mapstructure:"updater_interval" yaml:"updater_interval"`
	UPSInterval                   time.Duration     `mapstructure:"ups_interval" yaml:"ups_interval"`
	SyslogTagWireless             string            `mapstructure:"syslog_tag_wireless" yaml:"syslog_tag_wireless"`
	SyslogTagL2TP                 string            `mapstructure:"syslog_tag_l2tp" yaml:"syslog_tag_l2tp"`
	MacAddressMapping             map[string]string `mapstructure:"mac_address_mapping" yaml:"mac_address_mapping"`
	IgnoreUnknownMacAddress       bool              `mapstructure:"ignore_unknown_mac_address" yaml:"ignore_unknown_mac_address"`
	TopicInterfaceConnect         mqtt.Topic        `mapstructure:"topic_interface_connect" yaml:"topic_interface_connect"`
	TopicInterfaceLastConnect     mqtt.Topic        `mapstructure:"topic_interface_last_connect" yaml:"topic_interface_last_connect"`
	TopicInterfaceLastDisconnect  mqtt.Topic        `mapstructure:"topic_interface_last_disconnect" yaml:"topic_interface_last_disconnect"`
	TopicPackagesInstalledVersion mqtt.Topic        `mapstructure:"topic_packages_installed_version" yaml:"topic_packages_installed_version"`
	TopicPackagesLatestVersion    mqtt.Topic        `mapstructure:"topic_packages_latest_version" yaml:"topic_packages_latest_version"`
	TopicFirmwareInstalledVersion mqtt.Topic        `mapstructure:"topic_firmware_installed_version" yaml:"topic_firmware_installed_version"`
	TopicFirmwareLatestVersion    mqtt.Topic        `mapstructure:"topic_firmware_latest_version" yaml:"topic_firmware_latest_version"`
	TopicSyslog                   mqtt.Topic        `mapstructure:"topic_syslog" yaml:"topic_syslog"`
	TopicUPSBatteryCharge         mqtt.Topic        `mapstructure:"topic_ups_battery_charge" yaml:"topic_ups_battery_charge"`
	TopicUPSBatteryVoltage        mqtt.Topic        `mapstructure:"topic_ups_battery_voltage" yaml:"topic_ups_battery_voltage"`
	TopicUPSBatteryRuntime        mqtt.Topic        `mapstructure:"topic_ups_battery_runtime" yaml:"topic_ups_battery_runtime"`
	TopicUPSInputVoltage          mqtt.Topic        `mapstructure:"topic_ups_input_voltage" yaml:"topic_ups_input_voltage"`
	TopicUPSLoad                  mqtt.Topic        `mapstructure:"topic_ups_load" yaml:"topic_ups_load"`
	TopicUPSStatus                mqtt.Topic        `mapstructure:"topic_ups_status" yaml:"topic_ups_status"`
}

func (t Type) ConfigDefaults() interface{} {
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
		UPSInterval:                   time.Minute,
		SyslogTagWireless:             "wifi",
		SyslogTagL2TP:                 "vpn",
		IgnoreUnknownMacAddress:       true,
		TopicInterfaceConnect:         prefix + "interface/+/+/+",
		TopicInterfaceLastConnect:     prefix + "connect/+/+",
		TopicInterfaceLastDisconnect:  prefix + "disconnect/+/+",
		TopicPackagesInstalledVersion: prefix + "packages/installed-version",
		TopicPackagesLatestVersion:    prefix + "packages/latest-version",
		TopicFirmwareInstalledVersion: prefix + "firmware/installed-version",
		TopicFirmwareLatestVersion:    prefix + "firmware/latest-version",
		TopicUPSBatteryCharge:         prefix + "ups/+/battery-charge",
		TopicUPSBatteryVoltage:        prefix + "ups/+/battery-voltage",
		TopicUPSBatteryRuntime:        prefix + "ups/+/battery-runtime",
		TopicUPSInputVoltage:          prefix + "ups/+/input-voltage",
		TopicUPSLoad:                  prefix + "ups/+/load",
		TopicUPSStatus:                prefix + "ups/+/status",
	}
}

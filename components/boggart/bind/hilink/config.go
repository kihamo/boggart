package hilink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address                   boggart.URL `valid:",required"`
	Debug                     bool
	LivenessInterval          time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout           time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	BalanceUpdaterInterval    time.Duration `mapstructure:"balance_interval" yaml:"balance_interval"`
	BalanceUpdaterTimeout     time.Duration `mapstructure:"balance_timeout" yaml:"balance_timeout"`
	SMSCheckerInterval        time.Duration `mapstructure:"sms_checker_interval" yaml:"sms_checker_interval"`
	SMSCheckerTimeout         time.Duration `mapstructure:"sms_checker_timeout" yaml:"sms_checker_timeout"`
	SystemUpdaterInterval     time.Duration `mapstructure:"system_interval" yaml:"system_interval"`
	SystemUpdaterTimeout      time.Duration `mapstructure:"system_timeout" yaml:"system_timeout"`
	SMSCommandsEnabled        bool          `mapstructure:"sms_commands_enabled" yaml:"sms_commands_enabled"`
	SMSCommandsAllowedPhones  []string      `mapstructure:"sms_commands_allowed_phones" yaml:"sms_commands_allowed_phones"`
	TopicUSSDSend             mqtt.Topic    `mapstructure:"topic_ussd_send" yaml:"topic_ussd_send"`
	TopicUSSDResult           mqtt.Topic    `mapstructure:"topic_ussd_result" yaml:"topic_ussd_result"`
	TopicReboot               mqtt.Topic    `mapstructure:"topic_reboot" yaml:"topic_reboot"`
	TopicSMS                  mqtt.Topic    `mapstructure:"topic_sms" yaml:"topic_sms"`
	TopicBalance              mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicOperator             mqtt.Topic    `mapstructure:"topic_operator" yaml:"topic_operator"`
	TopicLimitInternetTraffic mqtt.Topic    `mapstructure:"topic_limits_internet_traffic" yaml:"topic_limits_internet_traffic"`
	TopicSignalRSSI           mqtt.Topic    `mapstructure:"topic_signal_rssi" yaml:"topic_signal_rssi"`
	TopicSignalRSRP           mqtt.Topic    `mapstructure:"topic_signal_rsrp" yaml:"topic_signal_rsrp"`
	TopicSignalRSRQ           mqtt.Topic    `mapstructure:"topic_signal_rsrq" yaml:"topic_signal_rsrq"`
	TopicSignalSINR           mqtt.Topic    `mapstructure:"topic_signal_sinr" yaml:"topic_signal_sinr"`
	TopicSignalLevel          mqtt.Topic    `mapstructure:"topic_signal_level" yaml:"topic_signal_level"`
	TopicConnectionTime       mqtt.Topic    `mapstructure:"topic_connection_time" yaml:"topic_connection_time"`
	TopicConnectionDownload   mqtt.Topic    `mapstructure:"topic_connection_download" yaml:"topic_connection_download"`
	TopicConnectionUpload     mqtt.Topic    `mapstructure:"topic_connection_upload" yaml:"topic_connection_upload"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/hilink/+/"

	return &Config{
		LivenessInterval:          time.Minute,
		LivenessTimeout:           time.Second * 5,
		BalanceUpdaterInterval:    time.Hour,
		BalanceUpdaterTimeout:     time.Second * 30,
		SMSCheckerInterval:        time.Minute,
		SMSCheckerTimeout:         time.Second * 30,
		SystemUpdaterInterval:     time.Minute,
		SystemUpdaterTimeout:      time.Second * 30,
		TopicUSSDSend:             prefix + "ussd/send",
		TopicUSSDResult:           prefix + "ussd",
		TopicReboot:               prefix + "reboot",
		TopicSMS:                  prefix + "sms",
		TopicBalance:              prefix + "balance",
		TopicOperator:             prefix + "operator",
		TopicLimitInternetTraffic: prefix + "limits/internet-traffic",
		TopicSignalRSSI:           prefix + "signal/rssi",
		TopicSignalRSRP:           prefix + "signal/rsrp",
		TopicSignalRSRQ:           prefix + "signal/rsrq",
		TopicSignalSINR:           prefix + "signal/sinr",
		TopicSignalLevel:          prefix + "signal/level",
		TopicConnectionTime:       prefix + "connection/time",
		TopicConnectionDownload:   prefix + "connection/download",
		TopicConnectionUpload:     prefix + "connection/upload",
	}
}

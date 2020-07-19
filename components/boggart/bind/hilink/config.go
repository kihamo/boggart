package hilink

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

	Address                     types.URL `valid:",required"`
	Debug                       bool
	SMSCommandsEnabled          bool          `mapstructure:"sms_commands_enabled" yaml:"sms_commands_enabled"`
	CleanerSpecial              bool          `mapstructure:"cleaner_special" yaml:"cleaner_special"`
	BalanceUpdaterInterval      time.Duration `mapstructure:"balance_interval" yaml:"balance_interval"`
	BalanceUpdaterTimeout       time.Duration `mapstructure:"balance_timeout" yaml:"balance_timeout"`
	LimitTrafficUpdaterInterval time.Duration `mapstructure:"limit_traffic_interval" yaml:"limit_traffic_interval"`
	LimitTrafficUpdaterTimeout  time.Duration `mapstructure:"limit_traffic_timeout" yaml:"limit_traffic_timeout"`
	SMSCheckerInterval          time.Duration `mapstructure:"sms_checker_interval" yaml:"sms_checker_interval"`
	SMSCheckerTimeout           time.Duration `mapstructure:"sms_checker_timeout" yaml:"sms_checker_timeout"`
	SystemUpdaterInterval       time.Duration `mapstructure:"system_interval" yaml:"system_interval"`
	SystemUpdaterTimeout        time.Duration `mapstructure:"system_timeout" yaml:"system_timeout"`
	CleanerInterval             time.Duration `mapstructure:"cleaner_interval" yaml:"cleaner_interval"`
	CleanerDuration             time.Duration `mapstructure:"cleaner_duration" yaml:"cleaner_duration"`
	SMSCommandsAllowedPhones    []string      `mapstructure:"sms_commands_allowed_phones" yaml:"sms_commands_allowed_phones"`
	TopicUSSDSend               mqtt.Topic    `mapstructure:"topic_ussd_send" yaml:"topic_ussd_send"`
	TopicUSSDResult             mqtt.Topic    `mapstructure:"topic_ussd_result" yaml:"topic_ussd_result"`
	TopicReboot                 mqtt.Topic    `mapstructure:"topic_reboot" yaml:"topic_reboot"`
	TopicSMS                    mqtt.Topic    `mapstructure:"topic_sms" yaml:"topic_sms"`
	TopicSMSUnread              mqtt.Topic    `mapstructure:"topic_sms_unread" yaml:"topic_sms_unread"`
	TopicSMSInbox               mqtt.Topic    `mapstructure:"topic_sms_inbox" yaml:"topic_sms_inbox"`
	TopicBalance                mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicOperator               mqtt.Topic    `mapstructure:"topic_operator" yaml:"topic_operator"`
	TopicLimitInternetTraffic   mqtt.Topic    `mapstructure:"topic_limits_internet_traffic" yaml:"topic_limits_internet_traffic"`
	TopicSignalRSSI             mqtt.Topic    `mapstructure:"topic_signal_rssi" yaml:"topic_signal_rssi"`
	TopicSignalRSRP             mqtt.Topic    `mapstructure:"topic_signal_rsrp" yaml:"topic_signal_rsrp"`
	TopicSignalRSRQ             mqtt.Topic    `mapstructure:"topic_signal_rsrq" yaml:"topic_signal_rsrq"`
	TopicSignalSINR             mqtt.Topic    `mapstructure:"topic_signal_sinr" yaml:"topic_signal_sinr"`
	TopicSignalLevel            mqtt.Topic    `mapstructure:"topic_signal_level" yaml:"topic_signal_level"`
	TopicConnectionTime         mqtt.Topic    `mapstructure:"topic_connection_time" yaml:"topic_connection_time"`
	TopicConnectionDownload     mqtt.Topic    `mapstructure:"topic_connection_download" yaml:"topic_connection_download"`
	TopicConnectionUpload       mqtt.Topic    `mapstructure:"topic_connection_upload" yaml:"topic_connection_upload"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/hilink/+/"

	return &Config{
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		BalanceUpdaterInterval:      time.Hour,
		BalanceUpdaterTimeout:       time.Second * 30,
		LimitTrafficUpdaterInterval: time.Hour,
		LimitTrafficUpdaterTimeout:  time.Second * 30,
		SMSCheckerInterval:          time.Minute,
		SMSCheckerTimeout:           time.Second * 30,
		SystemUpdaterInterval:       time.Minute,
		SystemUpdaterTimeout:        time.Second * 30,
		CleanerInterval:             time.Hour,
		CleanerSpecial:              true,
		CleanerDuration:             time.Hour * 24 * 90,
		TopicUSSDSend:               prefix + "ussd/send",
		TopicUSSDResult:             prefix + "ussd",
		TopicReboot:                 prefix + "reboot",
		TopicSMS:                    prefix + "sms",
		TopicSMSUnread:              prefix + "sms/unread",
		TopicSMSInbox:               prefix + "sms/inbox",
		TopicBalance:                prefix + "balance",
		TopicOperator:               prefix + "operator",
		TopicLimitInternetTraffic:   prefix + "limits/internet-traffic",
		TopicSignalRSSI:             prefix + "signal/rssi",
		TopicSignalRSRP:             prefix + "signal/rsrp",
		TopicSignalRSRQ:             prefix + "signal/rsrq",
		TopicSignalSINR:             prefix + "signal/sinr",
		TopicSignalLevel:            prefix + "signal/level",
		TopicConnectionTime:         prefix + "connection/time",
		TopicConnectionDownload:     prefix + "connection/download",
		TopicConnectionUpload:       prefix + "connection/upload",
	}
}

package smcenter

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

	ProviderLink            types.URL `mapstructure:"provider_link" yaml:"provider_link"  valid:"required"`
	Accounts                []string
	Phone                   string `valid:"required"`
	Password                string `valid:"required"`
	ProviderBaseURL         string `mapstructure:"provider_base_url" yaml:"provider_base_url"`
	ProviderBillContentType string `mapstructure:"provider_bill_content_type" yaml:"provider_bill_content_type"`
	AutoRegisterIfNotExists bool   `mapstructure:"auto_register_if_not_exists" yaml:"auto_register_if_not_exists"`
	Debug                   bool
	UpdaterInterval         time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicAccountBalance     mqtt.Topic    `mapstructure:"topic_account_balance" yaml:"topic_account_balance"`
	TopicAccountBill        mqtt.Topic    `mapstructure:"topic_account_bill" yaml:"topic_account_bill"`
	TopicMeterCheckupDate   mqtt.Topic    `mapstructure:"topic_meter_checkup_date" yaml:"topic_meter_checkup_date"`
	TopicMeterValue         mqtt.Topic    `mapstructure:"topic_meter_value" yaml:"topic_meter_value"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/smcenter/+/"

	cfg := &Config{
		ProbesConfig:            di.ProbesConfigDefaults(),
		LoggerConfig:            di.LoggerConfigDefaults(),
		AutoRegisterIfNotExists: true,
		Debug:                   false,
		UpdaterInterval:         time.Hour,
		TopicAccountBalance:     prefix + "balance",
		TopicAccountBill:        prefix + "bill",
		TopicMeterValue:         prefix + "meter/+/value",
		TopicMeterCheckupDate:   prefix + "meter/+/checkup",
		ProviderBaseURL:         t.BaseURL,
		ProviderBillContentType: t.BillContentType,
	}

	if t.Link != nil {
		cfg.ProviderLink.URL = *t.Link
	}

	return cfg
}

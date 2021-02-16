package dom24

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Accounts                []string
	Phone                   string `valid:"required"`
	Password                string `valid:"required"`
	AutoRegisterIfNotExists bool   `mapstructure:"auto_register_if_not_exists" yaml:"auto_register_if_not_exists"`
	Debug                   bool
	UpdaterInterval         time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	TopicAccountBalance     mqtt.Topic    `mapstructure:"topic_account_balance" yaml:"topic_account_balance"`
	TopicAccountBill        mqtt.Topic    `mapstructure:"topic_account_bill" yaml:"topic_account_bill"`
}

func (Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/dom24/"

	return &Config{
		ProbesConfig:            di.ProbesConfigDefaults(),
		LoggerConfig:            di.LoggerConfigDefaults(),
		AutoRegisterIfNotExists: true,
		Debug:                   false,
		UpdaterInterval:         time.Hour,
		TopicAccountBalance:     prefix + "+/balance",
		TopicAccountBill:        prefix + "+/bill",
	}
}

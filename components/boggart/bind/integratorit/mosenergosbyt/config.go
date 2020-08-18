package mosenergosbyt

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug                  bool
	Login                  string `valid:"required"`
	Password               string `valid:"required"`
	Account                string
	UpdaterInterval        time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	BalanceDetailsInterval time.Duration `mapstructure:"balance_details_interval" yaml:"balance_details_interval"`
	TopicBalance           mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicServiceBalance    mqtt.Topic    `mapstructure:"topic_service_balance" yaml:"topic_service_balance"`
	TopicLastBill          mqtt.Topic    `mapstructure:"topic_last_bill" yaml:"topic_last_bill"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/mosenergosbyt/+/"

	return &Config{
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		UpdaterInterval:        time.Hour,
		BalanceDetailsInterval: time.Hour * 24 * 31 * 2, // нужно минимум 2 месяца, что бы счет попал в выборку
		TopicBalance:           prefix + "balance",
		TopicServiceBalance:    prefix + "balance/+",
		TopicLastBill:          prefix + "bill/last",
	}
}

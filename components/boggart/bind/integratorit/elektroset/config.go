package elektroset

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
	Login                  string        `valid:"required"`
	Password               string        `valid:"required"`
	UpdaterInterval        time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	BalanceDetailsInterval time.Duration `mapstructure:"balance_details_interval" yaml:"balance_details_interval"`
	HouseID                uint64        `mapstructure:"house_id" yaml:"house_id"`
	TopicBalance           mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicServiceBalance    mqtt.Topic    `mapstructure:"topic_service_balance" yaml:"topic_service_balance"`
	TopicMeterValue        mqtt.Topic    `mapstructure:"topic_meter_value" yaml:"topic_meter_value"`
	TopicMeterDate         mqtt.Topic    `mapstructure:"topic_meter_date" yaml:"topic_meter_date"`
	TopicMeterCheckupDate  mqtt.Topic    `mapstructure:"topic_meter_checkup_date" yaml:"topic_meter_checkup_date"`
	TopicLastBill          mqtt.Topic    `mapstructure:"topic_last_bill" yaml:"topic_last_bill"`
}

func (Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/elektroset/+/"

	return &Config{
		ProbesConfig:           di.ProbesConfigDefaults(),
		LoggerConfig:           di.LoggerConfigDefaults(),
		UpdaterInterval:        time.Hour,
		BalanceDetailsInterval: time.Hour * 24 * 31 * 2, // нужно минимум 2 месяца, что бы счет попал в выборку
		TopicBalance:           prefix + "balance",
		TopicServiceBalance:    prefix + "service/+/balance",
		TopicMeterValue:        prefix + "meter/+/+/value",
		TopicMeterDate:         prefix + "meter/+/+/date",
		TopicMeterCheckupDate:  prefix + "meter/+/checkup",
		TopicLastBill:          prefix + "bill/last",
	}
}

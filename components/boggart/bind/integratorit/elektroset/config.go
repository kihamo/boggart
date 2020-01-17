package elektroset

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`

	Debug                  bool
	Login                  string        `valid:"required"`
	Password               string        `valid:"required"`
	UpdaterInterval        time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	BalanceDetailsInterval time.Duration `mapstructure:"balance_details_interval" yaml:"balance_details_interval"`
	TopicBalance           mqtt.Topic    `mapstructure:"topic_balance" yaml:"topic_balance"`
	TopicServiceBalance    mqtt.Topic    `mapstructure:"topic_service_balance" yaml:"topic_service_balance"`
	TopicMeterValueT1      mqtt.Topic    `mapstructure:"topic_meter_value_t1" yaml:"topic_meter_value_t1"`
	TopicMeterValueT2      mqtt.Topic    `mapstructure:"topic_meter_value_t2" yaml:"topic_meter_value_t2"`
	TopicMeterValueT3      mqtt.Topic    `mapstructure:"topic_meter_value_t3" yaml:"topic_meter_value_t3"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/elektroset/+/"

	return &Config{
		UpdaterInterval:        time.Hour,
		BalanceDetailsInterval: time.Hour * 24 * 31,
		TopicBalance:           prefix + "balance",
		TopicServiceBalance:    prefix + "balance/+",
		TopicMeterValueT1:      prefix + "meter/1",
		TopicMeterValueT2:      prefix + "meter/2",
		TopicMeterValueT3:      prefix + "meter/3",
	}
}

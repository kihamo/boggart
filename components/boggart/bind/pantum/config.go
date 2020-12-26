package pantum

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

	Address              types.URL `valid:",required"`
	Debug                bool
	UpdaterInterval      time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout       time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	TopicProductID       mqtt.Topic    `mapstructure:"topic_product_id" yaml:"topic_product_id"`
	TopicTonerRemain     mqtt.Topic    `mapstructure:"topic_toner_remain" yaml:"topic_toner_remain"`
	TopicPrinterStatus   mqtt.Topic    `mapstructure:"topic_printer_status" yaml:"topic_printer_status"`
	TopicCartridgeStatus mqtt.Topic    `mapstructure:"topic_cartridge_status" yaml:"topic_cartridge_status"`
	TopicDrumStatus      mqtt.Topic    `mapstructure:"topic_drum_status" yaml:"topic_drum_status"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/pantum/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 30
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig:         probesConfig,
		LoggerConfig:         di.LoggerConfigDefaults(),
		UpdaterInterval:      time.Minute,
		UpdaterTimeout:       time.Second * 10,
		TopicProductID:       prefix + "state/product-id",
		TopicTonerRemain:     prefix + "state/toner",
		TopicPrinterStatus:   prefix + "state/printer",
		TopicCartridgeStatus: prefix + "state/cartridge",
		TopicDrumStatus:      prefix + "state/drum",
	}
}

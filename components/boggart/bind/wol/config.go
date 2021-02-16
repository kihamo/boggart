package wol

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	TopicWOL                mqtt.Topic `mapstructure:"topic_wol" yaml:"topic_wol"`
	TopicWOLWithIPAndSubnet mqtt.Topic `mapstructure:"topic_wol_with_ip_and_subnet" yaml:"topic_wol_with_ip_and_subnet"`
}

func (Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/wol/"

	return &Config{
		TopicWOL:                prefix + "+",
		TopicWOLWithIPAndSubnet: prefix + "+/+/+",
	}
}

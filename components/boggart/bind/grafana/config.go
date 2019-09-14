package grafana

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Address         boggart.URL `valid:",required"`
	Name            string      `valid:",required"`
	ApiKey          string      `mapstructure:"api_key" yaml:"api_key"`
	Username        string
	Password        string
	Dashboards      []int64
	TopicAnnotation mqtt.Topic `mapstructure:"topic_annotation" yaml:"topic_annotation"`
}

func (t Type) Config() interface{} {
	return &Config{
		TopicAnnotation: boggart.ComponentName + "/grafana/+/annotation",
	}
}

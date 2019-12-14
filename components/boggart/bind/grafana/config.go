package grafana

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	boggart.BindConfig `mapstructure:",squash" yaml:",inline"`

	Address         boggart.URL `valid:",required"`
	Debug           bool
	Name            string `valid:",required"`
	Dashboards      []int64
	TopicAnnotation mqtt.Topic `mapstructure:"topic_annotation" yaml:"topic_annotation"`
}

func (t Type) Config() interface{} {
	return &Config{
		BindConfig: boggart.BindConfig{
			ReadinessPeriod:  time.Minute,
			ReadinessTimeout: time.Second * 5,
		},
		TopicAnnotation: boggart.ComponentName + "/grafana/+/annotation",
	}
}

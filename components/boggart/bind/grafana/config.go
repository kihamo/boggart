package grafana

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

	Address         types.URL `valid:",required"`
	Debug           bool
	Dashboards      []int64
	TopicAnnotation mqtt.Topic `mapstructure:"topic_annotation" yaml:"topic_annotation"`
}

func (t Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute
	probesConfig.ReadinessTimeout = time.Second * 5

	return &Config{
		ProbesConfig:    probesConfig,
		LoggerConfig:    di.LoggerConfigDefaults(),
		TopicAnnotation: boggart.ComponentName + "/grafana/+/annotation",
	}
}

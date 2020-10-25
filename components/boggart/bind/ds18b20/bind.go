package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/yryz/ds18b20"
)

type Bind struct {
	di.LoggerBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	config *Config
}

func (b *Bind) Temperature() (float64, error) {
	return ds18b20.Temperature(b.config.Address)
}

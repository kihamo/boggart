package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/yryz/ds18b20"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config *Config
}

func (b *Bind) Temperature() (float64, error) {
	return ds18b20.Temperature(b.SerialNumber())
}

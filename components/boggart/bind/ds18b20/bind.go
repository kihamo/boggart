package ds18b20

import (
	"time"

	"github.com/yryz/ds18b20"

	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	livenessInterval time.Duration
	livenessTimeout  time.Duration
	updaterInterval  time.Duration
}

func (b *Bind) Temperature() (float64, error) {
	return ds18b20.Temperature(b.SerialNumber())
}

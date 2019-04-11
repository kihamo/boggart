package ds18b20

import (
	"time"

	"github.com/yryz/ds18b20"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	temperature *atomic.Float32Null

	livenessInterval time.Duration
	livenessTimeout  time.Duration
	updaterInterval  time.Duration
}

func (b *Bind) Temperature() (float64, error) {
	return ds18b20.Temperature(b.SerialNumber())
}

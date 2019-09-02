package pulsar

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
)

const (
	InputScale = 1000
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config   *Config
	provider *pulsar.HeatMeter

	updaterInterval time.Duration
}

func (b *Bind) inputVolume(pulses float32, offset float32) float32 {
	return (offset*InputScale + pulses*10) / InputScale
}

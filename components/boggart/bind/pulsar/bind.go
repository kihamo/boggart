package pulsar

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
)

const (
	InputScale = 1000
)

type Bind struct {
	temperatureIn    *atomic.Float32
	temperatureOut   *atomic.Float32
	temperatureDelta *atomic.Float32
	energy           *atomic.Float32
	consumption      *atomic.Float32
	capacity         *atomic.Float32
	power            *atomic.Float32
	input1           *atomic.Float32
	input2           *atomic.Float32
	input3           *atomic.Float32
	input4           *atomic.Float32

	boggart.BindBase
	boggart.BindMQTT

	config   *Config
	provider *pulsar.HeatMeter

	updaterInterval time.Duration
}

func (b *Bind) inputVolume(pulses float32, offset float32) float32 {
	return (offset*InputScale + pulses*10) / InputScale
}

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
	temperatureIn    *atomic.Float32Null
	temperatureOut   *atomic.Float32Null
	temperatureDelta *atomic.Float32Null
	energy           *atomic.Float32Null
	consumption      *atomic.Float32Null
	capacity         *atomic.Float32Null
	power            *atomic.Float32Null
	input1           *atomic.Float32Null
	input2           *atomic.Float32Null
	input3           *atomic.Float32Null
	input4           *atomic.Float32Null

	boggart.BindBase
	boggart.BindMQTT

	config   *Config
	provider *pulsar.HeatMeter

	updaterInterval time.Duration
}

func (b *Bind) inputVolume(pulses float32, offset float32) float32 {
	return (offset*InputScale + pulses*10) / InputScale
}

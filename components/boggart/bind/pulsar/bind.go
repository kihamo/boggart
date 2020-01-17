package pulsar

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/pulsar"
)

const (
	InputScale = 1000
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.ProbesBind

	config   *Config
	provider *pulsar.HeatMeter
	address  string
}

func (b *Bind) inputVolume(pulses float32, offset float32) float32 {
	return (offset*InputScale + pulses*10) / InputScale
}

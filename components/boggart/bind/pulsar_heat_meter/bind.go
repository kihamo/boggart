package pulsar_heat_meter

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
)

const (
	InputScale = 1000
)

type Bind struct {
	temperatureIn    uint64
	temperatureOut   uint64
	temperatureDelta uint64
	energy           uint64
	consumption      uint64
	capacity         uint64
	power            uint64
	input1           uint64
	input2           uint64
	input3           uint64
	input4           uint64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	config   *Config
	provider *pulsar.HeatMeter

	updaterInterval time.Duration
}

func (b *Bind) inputVolume(pulses uint64, offset float64) float64 {
	return (offset*InputScale + float64(pulses*10)) / InputScale
}

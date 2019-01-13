package broadlink

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricSP3SPower = snitch.NewGauge(boggart.ComponentName+"_bind_broadlink_sp3s_power_watt", "Broadlink SP3S socket current power")
)

func (b *BindSP3S) Describe(ch chan<- *snitch.Description) {
	serialNumber := b.SerialNumber()

	metricSP3SPower.With("serial_number", serialNumber).Describe(ch)
}

func (b *BindSP3S) Collect(ch chan<- snitch.Metric) {
	serialNumber := b.SerialNumber()

	metricSP3SPower.With("serial_number", serialNumber).Collect(ch)
}

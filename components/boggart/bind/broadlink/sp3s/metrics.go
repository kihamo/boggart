package sp3s

import (
	"github.com/kihamo/snitch"
)

var (
	metricPower = snitch.NewGauge("power_watt", "Current power")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricPower.With("serial_number", b.config().MAC.String()).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricPower.With("serial_number", b.config().MAC.String()).Collect(ch)
}

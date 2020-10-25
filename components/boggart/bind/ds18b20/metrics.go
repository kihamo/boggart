package ds18b20

import (
	"github.com/kihamo/snitch"
)

var (
	metricValue = snitch.NewGauge("value_celsius", "Value in celsius")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricValue.With("serial_number", b.config.Address).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricValue.With("serial_number", b.config.Address).Collect(ch)
}

package neptun

import (
	"github.com/kihamo/snitch"
)

var (
	metricCounterValue = snitch.NewGauge("counter_value", "Counter value in cubic meters")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricCounterValue.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricCounterValue.Collect(ch)
}

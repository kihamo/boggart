package ping

import (
	"github.com/kihamo/snitch"
)

var (
	metricLatency = snitch.NewGauge("latency_milliseconds", "The network ping latency in milliseconds")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricLatency.With("host", b.config.Hostname).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricLatency.With("host", b.config.Hostname).Collect(ch)
}

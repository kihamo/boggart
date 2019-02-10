package service

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricLatency = snitch.NewGauge(boggart.ComponentName+"_bind_network_service_latency_milliseconds", "The network ping latency in milliseconds")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricLatency.With("address", b.address).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricLatency.With("address", b.address).Collect(ch)
}

package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricProbes = snitch.NewCounter(boggart.ComponentName+"_probes", "Readiness check")
)

func (c *Component) Describe(ch chan<- *snitch.Description) {
	metricProbes.Describe(ch)

	c.binds.Range(func(_ interface{}, item interface{}) bool {
		if collector, ok := item.(*BindItem).Bind().(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (c *Component) Collect(ch chan<- snitch.Metric) {
	metricProbes.Collect(ch)

	c.binds.Range(func(_ interface{}, item interface{}) bool {
		if collector, ok := item.(*BindItem).Bind().(snitch.Collector); ok {
			collector.Collect(ch)
		}

		return true
	})
}

func (c *Component) Metrics() snitch.Collector {
	<-c.application.ReadyComponent(c.Name())
	return c
}

package internal

import (
	"github.com/kihamo/snitch"
)

type MetricsCollector struct {
	component *Component
}

func NewMetricsCollector(component *Component) *MetricsCollector {
	return &MetricsCollector{
		component: component,
	}
}

func (c *MetricsCollector) Describe(ch chan<- *snitch.Description) {
	c.component.devices.Describe(ch)

	c.DescribeMercury(ch)
	c.DescribePulsar(ch)
}

func (c *MetricsCollector) Collect(ch chan<- snitch.Metric) {
	c.component.devices.Collect(ch)

	c.CollectMercury(ch)
	c.CollectPulsar(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return c.collector
}

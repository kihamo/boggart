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
	c.DescribeDoors(ch)
	c.DescribeMikrotik(ch)
	c.DescribePulsar(ch)
	c.DescribeSoftVideo(ch)
}

func (c *MetricsCollector) Collect(ch chan<- snitch.Metric) {
	c.CollectDoors(ch)
	c.CollectMikrotik(ch)
	c.CollectPulsar(ch)
	c.CollectSoftVideo(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return c.collector
}

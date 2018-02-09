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

	c.DescribeDoors(ch)
	c.DescribeMercury(ch)
	c.DescribeMikrotik(ch)
	c.DescribeMobile(ch)
	c.DescribePulsar(ch)
	c.DescribeSoftVideo(ch)
}

func (c *MetricsCollector) Collect(ch chan<- snitch.Metric) {
	c.component.devices.Collect(ch)

	c.CollectDoors(ch)
	c.CollectMercury(ch)
	c.CollectMikrotik(ch)
	c.CollectMobile(ch)
	c.CollectPulsar(ch)
	c.CollectSoftVideo(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return c.collector
}

package internal

import (
	"github.com/kihamo/snitch"
)

func (c *Component) Describe(ch chan<- *snitch.Description) {
	c.binds.Range(func(_ interface{}, item interface{}) bool {
		if collector, ok := item.(*BindItem).Bind().(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (c *Component) Collect(ch chan<- snitch.Metric) {
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

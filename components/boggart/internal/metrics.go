package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/snitch"
)

var (
	metricProbes     = snitch.NewCounter(boggart.ComponentName+"_probes", "Readiness check")
	metricBindStatus = snitch.NewCounter(boggart.ComponentName+"_status_bind", "Bind status")
)

func (c *Component) Describe(ch chan<- *snitch.Description) {
	metricProbes.Describe(ch)
	metricBindStatus.Describe(ch)

	c.tasksManager.Describe(ch)
	c.binds.Range(func(_ interface{}, item interface{}) bool {
		if bindSupport, ok := di.MetricsContainerBind(item.(*BindItem).Bind()); ok {
			bindSupport.Describe(ch)
		}

		return true
	})
}

func (c *Component) Collect(ch chan<- snitch.Metric) {
	metricProbes.Collect(ch)
	metricBindStatus.Collect(ch)

	c.tasksManager.Collect(ch)
	c.binds.Range(func(_ interface{}, item interface{}) bool {
		if bindSupport, ok := di.MetricsContainerBind(item.(*BindItem).Bind()); ok {
			bindSupport.Collect(ch)
		}

		return true
	})
}

func (c *Component) Metrics() snitch.Collector {
	<-c.application.ReadyComponent(c.Name())
	return c
}

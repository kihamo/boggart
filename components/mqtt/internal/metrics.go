package internal

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/snitch"
)

const (
	MetricPublish   = mqtt.ComponentName + "_publish_total"
	MetricSubscribe = mqtt.ComponentName + "_subscribe_total"
)

var (
	metricPublish   = snitch.NewCounter(MetricPublish, "Total publish")
	metricSubscribe = snitch.NewCounter(MetricSubscribe, "Total subscribe")
)

type metricsCollector struct {
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricPublish.Describe(ch)
	metricSubscribe.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	metricPublish.Collect(ch)
	metricSubscribe.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{}
}

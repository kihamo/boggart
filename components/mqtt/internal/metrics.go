package internal

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/snitch"
)

const (
	MetricPublish         = mqtt.ComponentName + "_publish_total"
	MetricSubscribe       = mqtt.ComponentName + "_subscribe_total"
	MetricSubscriberCalls = mqtt.ComponentName + "_subscriber_calls_total"
	MetricConnect         = mqtt.ComponentName + "_connect_total"
	MetricConnectionLost  = mqtt.ComponentName + "_connection_lost_total"
)

var (
	metricPublish         = snitch.NewCounter(MetricPublish, "Total publish")
	metricSubscribe       = snitch.NewCounter(MetricSubscribe, "Total subscribe")
	metricSubscriberCalls = snitch.NewCounter(MetricSubscriberCalls, "Total subscriber calls")
	metricConnect         = snitch.NewCounter(MetricConnect, "Total connect")
	metricConnectionLost  = snitch.NewCounter(MetricConnectionLost, "Total connection lost")
)

type metricsCollector struct {
}

func (c *metricsCollector) Describe(ch chan<- *snitch.Description) {
	metricPublish.Describe(ch)
	metricSubscribe.Describe(ch)
	metricSubscriberCalls.Describe(ch)
	metricConnect.Describe(ch)
	metricConnectionLost.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- snitch.Metric) {
	metricPublish.Collect(ch)
	metricSubscribe.Collect(ch)
	metricSubscriberCalls.Collect(ch)
	metricConnect.Collect(ch)
	metricConnectionLost.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return &metricsCollector{}
}

package internal

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/snitch"
)

const (
	MetricPublish           = mqtt.ComponentName + "_publish_total"
	MetricSubscribe         = mqtt.ComponentName + "_subscribe_total"
	MetricSubscriberCalls   = mqtt.ComponentName + "_subscriber_calls_total"
	MetricConnect           = mqtt.ComponentName + "_connect_total"
	MetricConnectionLost    = mqtt.ComponentName + "_connection_lost_total"
	MetricPayloadCacheItems = mqtt.ComponentName + "_payload_cache_items_total"
	MetricPayloadCacheHit   = mqtt.ComponentName + "_payload_cache_hit_total"
	MetricPayloadCacheMiss  = mqtt.ComponentName + "_payload_cache_miss_total"
)

var (
	metricPublish           = snitch.NewCounter(MetricPublish, "Total publish")
	metricSubscribe         = snitch.NewCounter(MetricSubscribe, "Total subscribe")
	metricSubscriberCalls   = snitch.NewCounter(MetricSubscriberCalls, "Total subscriber calls")
	metricConnect           = snitch.NewCounter(MetricConnect, "Total connect")
	metricConnectionLost    = snitch.NewCounter(MetricConnectionLost, "Total connection lost")
	metricPayloadCacheItems = snitch.NewGauge(MetricPayloadCacheItems, "Count items of payload cache")
	metricPayloadCacheHit   = snitch.NewCounter(MetricPayloadCacheHit, "Hits of payload cache")
	metricPayloadCacheMiss  = snitch.NewCounter(MetricPayloadCacheMiss, "Misses of payload cache")
)

func (c *Component) Describe(ch chan<- *snitch.Description) {
	metricPublish.Describe(ch)
	metricSubscribe.Describe(ch)
	metricSubscriberCalls.Describe(ch)
	metricConnect.Describe(ch)
	metricConnectionLost.Describe(ch)
	metricPayloadCacheItems.Describe(ch)
	metricPayloadCacheHit.Describe(ch)
	metricPayloadCacheMiss.Describe(ch)
}

func (c *Component) Collect(ch chan<- snitch.Metric) {
	metricPayloadCacheItems.Set(float64(c.payloadCache.Len()))

	metricPublish.Collect(ch)
	metricSubscribe.Collect(ch)
	metricSubscriberCalls.Collect(ch)
	metricConnect.Collect(ch)
	metricConnectionLost.Collect(ch)
	metricPayloadCacheItems.Collect(ch)
	metricPayloadCacheHit.Collect(ch)
	metricPayloadCacheMiss.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return c
}

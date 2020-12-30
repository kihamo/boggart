package internal

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/snitch"
)

const (
	MetricPublish                   = mqtt.ComponentName + "_publish_total"
	MetricSubscribe                 = mqtt.ComponentName + "_subscribe_total"
	MetricSubscriberCalls           = mqtt.ComponentName + "_subscriber_calls_total"
	MetricConnect                   = mqtt.ComponentName + "_connect_total"
	MetricConnectionLost            = mqtt.ComponentName + "_connection_lost_total"
	MetricPayloadCacheItems         = mqtt.ComponentName + "_payload_cache_items_total"
	MetricPayloadCacheHit           = mqtt.ComponentName + "_payload_cache_hit_total"
	MetricPayloadCacheMiss          = mqtt.ComponentName + "_payload_cache_miss_total"
	MetricBrokerClientsConnected    = mqtt.ComponentName + "_broker_clients_connected"
	MetricBrokerClientsDisconnected = mqtt.ComponentName + "_broker_clients_disconnected"
	MetricBrokerClientsMaximum      = mqtt.ComponentName + "_broker_clients_maximum"
	MetricBrokerClientsTotal        = mqtt.ComponentName + "_broker_clients_total"
	MetricBrokerMessages            = mqtt.ComponentName + "_broker_messages"
	MetricBrokerSubscriptions       = mqtt.ComponentName + "_broker_subscriptions_total"
	MetricBrokerUpTime              = mqtt.ComponentName + "_broker_uptime_seconds"
	MetricBrokerVersion             = mqtt.ComponentName + "_broker_version"
	MetricBrokerBytes               = mqtt.ComponentName + "_broker_bytes"
	MetricBrokerHeapCurrent         = mqtt.ComponentName + "_broker_heap_current_bytes"
	MetricBrokerHeapMaximum         = mqtt.ComponentName + "_broker_heap_maximum_bytes"
	MetricBrokerLoadBytes           = mqtt.ComponentName + "_broker_load_bytes"
	MetricBrokerLoadConnections     = mqtt.ComponentName + "_broker_load_connections"
	MetricBrokerLoadMessages        = mqtt.ComponentName + "_broker_load_messages"
	MetricBrokerLoadPublish         = mqtt.ComponentName + "_broker_load_publish"
	MetricBrokerLoadSockets         = mqtt.ComponentName + "_broker_load_sockets"
	MetricBrokerMessagesStoredBytes = mqtt.ComponentName + "_broker_messages_stored_bytes"
	MetricBrokerPublishMessages     = mqtt.ComponentName + "_broker_publish_messages"
	MetricBrokerPublishBytes        = mqtt.ComponentName + "_broker_publish_bytes"
)

var (
	metricPublish                   = snitch.NewCounter(MetricPublish, "Total publish")
	metricSubscribe                 = snitch.NewCounter(MetricSubscribe, "Total subscribe")
	metricSubscriberCalls           = snitch.NewCounter(MetricSubscriberCalls, "Total subscriber calls")
	metricConnect                   = snitch.NewCounter(MetricConnect, "Total connect")
	metricConnectionLost            = snitch.NewCounter(MetricConnectionLost, "Total connection lost")
	metricPayloadCacheItems         = snitch.NewGauge(MetricPayloadCacheItems, "Count items of payload cache")
	metricPayloadCacheHit           = snitch.NewCounter(MetricPayloadCacheHit, "Hits of payload cache")
	metricPayloadCacheMiss          = snitch.NewCounter(MetricPayloadCacheMiss, "Misses of payload cache")
	metricBrokerClientsConnected    = snitch.NewGauge(MetricBrokerClientsConnected, "The number of currently connected client")
	metricBrokerClientsDisconnected = snitch.NewGauge(MetricBrokerClientsDisconnected, "The total number of persistent clients (with clean session disabled) that are registered at the broker but are currently disconnected")
	metricBrokerClientsMaximum      = snitch.NewGauge(MetricBrokerClientsMaximum, "The maximum number of active clients that have been connected to the broker. This is only calculated when the $SYS topic tree is updated, so short lived client connections may not be counted")
	metricBrokerClientsTotal        = snitch.NewGauge(MetricBrokerClientsTotal, "The total number of connected and disconnected clients with a persistent session currently connected and registered on the broker")
	metricBrokerMessages            = snitch.NewGauge(MetricBrokerMessages, "The total number of messages of any type since the broker started")
	metricBrokerSubscriptions       = snitch.NewGauge(MetricBrokerSubscriptions, "The total number of subscriptions active on the broker")
	metricBrokerUpTime              = snitch.NewGauge(MetricBrokerUpTime, "The amount of time in seconds the broker has been online")
	metricBrokerVersion             = snitch.NewGauge(MetricBrokerVersion, "The version of the broker")
	metricBrokerBytes               = snitch.NewGauge(MetricBrokerBytes, "The total number of bytes since the broker started")
	metricBrokerHeapCurrent         = snitch.NewGauge(MetricBrokerHeapCurrent, "The current size of the heap memory in use by mosquitto")
	metricBrokerHeapMaximum         = snitch.NewGauge(MetricBrokerHeapMaximum, "The largest amount of heap memory used by mosquitto")
	metricBrokerLoadBytes           = snitch.NewGauge(MetricBrokerLoadBytes, "The moving average of the number of bytes by the broker over different time interval")
	metricBrokerLoadConnections     = snitch.NewGauge(MetricBrokerLoadConnections, "The moving average of the number of CONNECT packets received by the broker over different time intervals")
	metricBrokerLoadMessages        = snitch.NewGauge(MetricBrokerLoadMessages, "The moving average of the number of all types of MQTT messages by the broker over different time intervals")
	metricBrokerLoadPublish         = snitch.NewGauge(MetricBrokerLoadPublish, "The moving average of the number of publish messages by the broker over different time intervals")
	metricBrokerLoadSockets         = snitch.NewGauge(MetricBrokerLoadSockets, "The moving average of the number of socket connections opened to the broker over different time intervals")
	metricBrokerMessagesStoredBytes = snitch.NewGauge(MetricBrokerMessagesStoredBytes, "The number of bytes currently held by message payloads in the message store. This includes retained messages and messages queued for durable clients")
	metricBrokerPublishMessages     = snitch.NewGauge(MetricBrokerPublishMessages, "The total number of publish messages")
	metricBrokerPublishBytes        = snitch.NewGauge(MetricBrokerPublishBytes, "The total number of PUBLISH bytes since the broker started")
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

	metricBrokerClientsConnected.Describe(ch)
	metricBrokerClientsDisconnected.Describe(ch)
	metricBrokerClientsMaximum.Describe(ch)
	metricBrokerClientsTotal.Describe(ch)
	metricBrokerMessages.Describe(ch)
	metricBrokerSubscriptions.Describe(ch)
	metricBrokerUpTime.Describe(ch)
	metricBrokerVersion.Describe(ch)
	metricBrokerBytes.Describe(ch)
	metricBrokerHeapCurrent.Describe(ch)
	metricBrokerHeapMaximum.Describe(ch)
	metricBrokerLoadBytes.Describe(ch)
	metricBrokerLoadConnections.Describe(ch)
	metricBrokerLoadMessages.Describe(ch)
	metricBrokerLoadPublish.Describe(ch)
	metricBrokerLoadSockets.Describe(ch)
	metricBrokerMessagesStoredBytes.Describe(ch)
	metricBrokerPublishMessages.Describe(ch)
	metricBrokerPublishBytes.Describe(ch)
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

	metricBrokerClientsConnected.Collect(ch)
	metricBrokerClientsDisconnected.Collect(ch)
	metricBrokerClientsMaximum.Collect(ch)
	metricBrokerClientsTotal.Collect(ch)
	metricBrokerMessages.Collect(ch)
	metricBrokerSubscriptions.Collect(ch)
	metricBrokerUpTime.Collect(ch)
	metricBrokerVersion.Collect(ch)
	metricBrokerBytes.Collect(ch)
	metricBrokerHeapCurrent.Collect(ch)
	metricBrokerHeapMaximum.Collect(ch)
	metricBrokerLoadBytes.Collect(ch)
	metricBrokerLoadConnections.Collect(ch)
	metricBrokerLoadMessages.Collect(ch)
	metricBrokerLoadPublish.Collect(ch)
	metricBrokerLoadSockets.Collect(ch)
	metricBrokerMessagesStoredBytes.Collect(ch)
	metricBrokerPublishMessages.Collect(ch)
	metricBrokerPublishBytes.Collect(ch)
}

func (c *Component) Metrics() snitch.Collector {
	return c
}

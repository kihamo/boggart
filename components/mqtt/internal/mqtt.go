package internal

// https://github.com/mqtt/mqtt.org/wiki/SYS-Topics
// https://mosquitto.org/man/mosquitto-8.html

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber("$SYS/broker/#", 0, c.callbackBroker),
	}
}

func (c *Component) callbackBroker(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	switch message.Topic() {
	// Required topics
	case "$SYS/broker/clients/connected", "$SYS/broker/clients/active":
		metricBrokerClientsConnected.Set(message.Float64())
	case "$SYS/broker/clients/disconnected", "$SYS/broker/clients/inactive":
		metricBrokerClientsDisconnected.Set(message.Float64())
	case "$SYS/broker/clients/maximum":
		metricBrokerClientsMaximum.Set(message.Float64())
	case "$SYS/broker/clients/total":
		metricBrokerClientsTotal.Set(message.Float64())
	case "$SYS/broker/messages/received":
		metricBrokerMessages.With("status", "received").Set(message.Float64())
	case "$SYS/broker/messages/sent":
		metricBrokerMessages.With("status", "sent").Set(message.Float64())
	case "$SYS/broker/subscriptions/count":
		metricBrokerSubscriptions.Set(message.Float64())
	case "$SYS/broker/uptime":
		if parts := strings.Fields(message.String()); len(parts) > 0 {
			if v, err := strconv.ParseFloat(parts[0], 64); err == nil {
				metricBrokerUpTime.Set(v)
			}
		}
	case "$SYS/broker/version":
		metricBrokerVersion.With("version", message.String()).Set(1)

		// Optional topics
	case "$SYS/broker/messages/inflight":
		metricBrokerMessages.With("status", "inflight").Set(message.Float64())
	case "$SYS/broker/messages/stored", "$SYS/broker/store/messages/count", "$SYS/broker/retained messages/count":
		metricBrokerMessages.With("status", "stored").Set(message.Float64())

		// Mosquitto topics
	case "$SYS/broker/bytes/received":
		metricBrokerBytes.With("status", "received").Set(message.Float64())
	case "$SYS/broker/bytes/sent":
		metricBrokerBytes.With("status", "sent").Set(message.Float64())
	case "$SYS/broker/heap/current":
		metricBrokerHeapCurrent.Set(message.Float64())
	case "$SYS/broker/heap/maximum":
		metricBrokerHeapMaximum.Set(message.Float64())
	case "$SYS/broker/store/messages/bytes":
		metricBrokerMessagesStoredBytes.Set(message.Float64())
	case "$SYS/broker/publish/messages/dropped":
		metricBrokerPublishMessages.With("status", "dropped").Set(message.Float64())
	case "$SYS/broker/publish/messages/received":
		metricBrokerPublishMessages.With("status", "received").Set(message.Float64())
	case "$SYS/broker/publish/messages/sent":
		metricBrokerPublishMessages.With("status", "sent").Set(message.Float64())
	case "$SYS/broker/publish/bytes/received":
		metricBrokerPublishBytes.With("status", "received").Set(message.Float64())
	case "$SYS/broker/publish/bytes/sent":
		metricBrokerPublishBytes.With("status", "sent").Set(message.Float64())
	default:
		topic := message.Topic().String()
		index := strings.LastIndex(topic, "/")
		if index < 1 {
			break
		}

		switch topic[:index] {
		case "$SYS/broker/load/connections":
			metricBrokerLoadConnections.With("interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/bytes/received":
			metricBrokerLoadBytes.With("status", "received", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/bytes/sent":
			metricBrokerLoadBytes.With("status", "sent", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/messages/received":
			metricBrokerLoadMessages.With("status", "received", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/messages/sent":
			metricBrokerLoadMessages.With("status", "sent", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/publish/dropped":
			metricBrokerLoadPublish.With("status", "dropped", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/publish/received":
			metricBrokerLoadPublish.With("status", "received", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/publish/sent":
			metricBrokerLoadPublish.With("status", "sent", "interval", topic[index+1:]).Set(message.Float64())
		case "$SYS/broker/load/sockets":
			metricBrokerLoadSockets.With("interval", topic[index+1:]).Set(message.Float64())

		default:
			c.logger.Warn("Unknown stats topic", "topic", topic)
		}
	}

	return nil
}

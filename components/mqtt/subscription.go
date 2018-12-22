package mqtt

import (
	"context"
	"strings"
	"sync/atomic"
)

type Subscription struct {
	calls     uint64
	topic     string
	qos       byte
	callbacks []MessageHandler
}

func NewSubscription(topic string, qos byte, callbacks ...MessageHandler) *Subscription {
	return &Subscription{
		topic:     topic,
		qos:       qos,
		callbacks: callbacks,
	}
}

func (c *Subscription) Calls() uint64 {
	return atomic.LoadUint64(&c.calls)
}

func (c *Subscription) Topic() string {
	return c.topic
}

func (c *Subscription) QOS() byte {
	return c.qos
}

func (c *Subscription) Merge(sub *Subscription) *Subscription {
	return &Subscription{
		calls:     atomic.LoadUint64(&c.calls),
		topic:     sub.topic,
		qos:       sub.qos,
		callbacks: append(c.callbacks, sub.callbacks...),
	}
}

func (c *Subscription) Callback(ctx context.Context, client Component, message Message) {
	atomic.AddUint64(&c.calls, 1)

	for _, callback := range c.callbacks {
		go func(c MessageHandler) {
			c(ctx, client, message)
		}(callback)
	}
}

func (c *Subscription) Len() int {
	return len(c.callbacks)
}

func (c *Subscription) Match(topic string) bool {
	return c.topic == topic || routeIncludesTopic(c.topic, topic)
}

func match(route []string, topic []string) bool {
	if len(route) == 0 {
		return len(topic) == 0
	}

	if len(topic) == 0 {
		return route[0] == "#"
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return match(route[1:], topic[1:])
	}

	return false
}

func routeIncludesTopic(route, topic string) bool {
	topic = strings.TrimRight(topic, "/")

	return match(RouteSplit(route), strings.Split(topic, "/"))
}

func RouteSplit(route string) []string {
	route = strings.TrimRight(route, "/")

	var result []string

	if strings.HasPrefix(route, "$share") {
		result = strings.Split(route, "/")[2:]
	} else {
		result = strings.Split(route, "/")
	}

	return result
}

package mqtt

import (
	"context"
)

type HasSubscribers interface {
	MQTTSubscribers() []Subscriber
}

type Subscriber interface {
	Topic() string
	QOS() byte
	Call(context.Context, Component, Message)
}

type SubscriberSimple struct {
	topic    string
	qos      byte
	callback MessageHandler
}

func NewSubscriber(topic string, qos byte, callback MessageHandler) *SubscriberSimple {
	return &SubscriberSimple{
		topic:    topic,
		qos:      qos,
		callback: callback,
	}
}

func (s *SubscriberSimple) Topic() string {
	return s.topic
}

func (s *SubscriberSimple) QOS() byte {
	return s.qos
}

func (s *SubscriberSimple) Call(ctx context.Context, client Component, message Message) {
	s.callback(ctx, client, message)
}

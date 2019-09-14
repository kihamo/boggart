package mqtt

import (
	"context"
)

type HasSubscribers interface {
	MQTTSubscribers() []Subscriber
}

type Subscriber interface {
	Topic() Topic
	QOS() byte
	Call(context.Context, Component, Message) error
}

type SubscriberSimple struct {
	topic    Topic
	qos      byte
	callback MessageHandler
}

func NewSubscriber(topic Topic, qos byte, callback MessageHandler) *SubscriberSimple {
	return &SubscriberSimple{
		topic:    topic,
		qos:      qos,
		callback: callback,
	}
}

func (s *SubscriberSimple) Topic() Topic {
	return s.topic
}

func (s *SubscriberSimple) QOS() byte {
	return s.qos
}

func (s *SubscriberSimple) Call(ctx context.Context, client Component, message Message) error {
	return s.callback(ctx, client, message)
}

package mqtt

import (
	"context"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Client() m.Client
	Publish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error
	PublishWithCache(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error
	PublishWithoutCache(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error
	PublishAsync(ctx context.Context, topic string, qos byte, retained bool, payload interface{})
	PublishAsyncWithCache(ctx context.Context, topic string, qos byte, retained bool, payload interface{})
	PublishAsyncWithoutCache(ctx context.Context, topic string, qos byte, retained bool, payload interface{})

	Unsubscribe(topic string) error
	UnsubscribeSubscriber(Subscriber) error
	UnsubscribeSubscribers([]Subscriber) error

	Subscribe(topic string, qos byte, callback MessageHandler) (Subscriber, error)
	SubscribeSubscriber(Subscriber) error

	Subscriptions() []*Subscription
}

type Message interface {
	m.Message

	UnmarshalJSON(interface{}) error
	IsTrue() bool
	IsFalse() bool
	Bool() bool
	String() string
}

type MessageHandler func(ctx context.Context, client Component, message Message) error

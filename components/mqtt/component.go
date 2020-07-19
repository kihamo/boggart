package mqtt

import (
	"context"
	"fmt"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Client() m.Client
	Publish(ctx context.Context, topic Topic, qos byte, retained bool, payload interface{}) error
	PublishWithCache(ctx context.Context, topic Topic, qos byte, retained bool, payload interface{}) error
	PublishWithoutCache(ctx context.Context, topic Topic, qos byte, retained bool, payload interface{}) error
	PublishAsync(ctx context.Context, topic Topic, qos byte, retained bool, payload interface{})
	PublishAsyncWithCache(ctx context.Context, topic Topic, qos byte, retained bool, payload interface{})
	PublishAsyncWithoutCache(ctx context.Context, topic Topic, qos byte, retained bool, payload interface{})

	Unsubscribe(topic Topic) error
	UnsubscribeSubscriber(Subscriber) error
	UnsubscribeSubscribers([]Subscriber) error

	Subscribe(topic Topic, qos byte, callback MessageHandler) (Subscriber, error)
	SubscribeSubscriber(Subscriber) error

	Subscriptions() []*Subscription

	OnConnectHandlerAdd(handler OnConnectHandler)
	CacheItems() []CacheItem
}

type Message interface {
	fmt.Stringer

	Duplicate() bool
	Qos() byte
	Retained() bool
	Topic() Topic
	MessageID() uint16
	Payload() []byte
	Ack()
	JSONUnmarshal(interface{}) error
	IsTrue() bool
	IsFalse() bool
	Bool() bool
	Base64() ([]byte, error)
	Len() int
}

type OnConnectHandler func(client Component, restore bool)
type MessageHandler func(ctx context.Context, client Component, message Message) error

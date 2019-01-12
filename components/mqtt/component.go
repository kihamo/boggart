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

	Unsubscribe(topic string) error
	UnsubscribeSubscriber(Subscriber) error
	UnsubscribeSubscribers([]Subscriber) error

	Subscribe(topic string, qos byte, callback MessageHandler) (Subscriber, error)
	SubscribeSubscriber(Subscriber) error

	Subscriptions() []*Subscription
}

type Message m.Message

type MessageHandler func(ctx context.Context, client Component, message Message)

package mqtt

import (
	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Client() m.Client
	Publish(topic string, qos byte, retained bool, payload interface{}) m.Token
	Subscribe(topic string, qos byte, callback m.MessageHandler) m.Token
	SubscribeMultiple(filters map[string]byte, callback m.MessageHandler) m.Token
	Unsubscribe(topics ...string) m.Token
	AddRoute(topic string, callback m.MessageHandler)
	OptionsReader() m.ClientOptionsReader
}

package mqtt

import (
	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Client() m.Client
	Publish(topic string, qos byte, retained bool, payload interface{}) m.Token
	Subscribe(subscriber Subscriber) m.Token
}

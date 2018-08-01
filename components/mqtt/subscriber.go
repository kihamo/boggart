package mqtt

import (
	m "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber interface {
	Filters() map[string]byte
	Callback(Component, m.Message)
}

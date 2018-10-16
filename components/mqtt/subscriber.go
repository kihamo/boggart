package mqtt

import (
	"context"

	m "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber interface {
	Filters() map[string]byte
	Callback(context.Context, Component, m.Message)
}

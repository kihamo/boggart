package subscribes

import (
	"context"
	"encoding/json"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/mmcloughlin/geohash"
)

type OwnTracksSubscribe struct{}

func NewOwnTracksSubscribe() *OwnTracksSubscribe {
	return &OwnTracksSubscribe{}
}

func (s *OwnTracksSubscribe) Filters() map[string]byte {
	return map[string]byte{
		"owntracks/#": 0,
	}
}

func (s *OwnTracksSubscribe) Callback(_ context.Context, client mqtt.Component, message m.Message) {
	var payload map[string]interface{}

	err := json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		return
	}

	t, ok := payload["_type"]
	if !ok || t != "location" {
		return
	}

	lat, ok := payload["lat"]
	if !ok {
		return
	}

	lon, ok := payload["lon"]
	if !ok {
		return
	}

	hash := geohash.Encode(lat.(float64), lon.(float64))
	client.Publish(message.Topic()+"/geohash", message.Qos(), message.Retained(), hash)
}

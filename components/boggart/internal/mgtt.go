package internal

import (
	"context"
	"encoding/json"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/mmcloughlin/geohash"
)

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	<-c.application.ReadyComponent(c.Name())

	if !c.config.Bool(boggart.ConfigOwnTracksEnabled) {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber("owntracks/#", 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
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
		}),
	}
}

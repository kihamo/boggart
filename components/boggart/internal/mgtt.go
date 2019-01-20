package internal

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/mmcloughlin/geohash"
)

const (
	MQTTSubscribeTopicOwnTracks         mqtt.Topic = "owntracks/+/+"
	MQTTPublishTopicOwnTracksGeoHash    mqtt.Topic = "owntracks/+/+/geohash"
	MQTTSubscribeTopicAnnotationGrafana mqtt.Topic = "annotation/grafana"
	MQTTSubscribeTopicMessenger         mqtt.Topic = "messenger/+/+"
)

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	<-c.application.ReadyComponent(c.Name())

	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicOwnTracks.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			if !c.config.Bool(boggart.ConfigMQTTOwnTracksEnabled) {
				return nil
			}

			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 2 {
				return errors.New("bad topic name")
			}

			var payload map[string]interface{}

			if err := json.Unmarshal(message.Payload(), &payload); err != nil {
				return err
			}

			t, ok := payload["_type"]
			if !ok || t != "location" {
				return errors.New("location not found in payload")
			}

			lat, ok := payload["lat"]
			if !ok {
				return errors.New("lat not found in payload")
			}

			lon, ok := payload["lon"]
			if !ok {
				return errors.New("lon not found in payload")
			}

			hash := geohash.Encode(lat.(float64), lon.(float64))
			return client.Publish(ctx, MQTTPublishTopicOwnTracksGeoHash.Format(route[len(route)-2], route[len(route)-1]), message.Qos(), message.Retained(), hash)
		}),
	}

	if c.application.HasComponent(annotations.ComponentName) {
		<-c.application.ReadyComponent(annotations.ComponentName)
		cmp := c.application.GetComponent(annotations.ComponentName).(annotations.Component)

		subscribers = append(subscribers, mqtt.NewSubscriber(MQTTSubscribeTopicAnnotationGrafana.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !c.config.Bool(boggart.ConfigMQTTAnnotationsEnabled) {
				return nil
			}

			var request struct {
				Title string   `json:"title,omitempty"`
				Text  string   `json:"text,omitempty"`
				Tags  []string `json:"tags,omitempty"`
			}

			if err := json.Unmarshal(message.Payload(), &request); err != nil {
				return err
			}

			return cmp.CreateInStorages(
				annotations.NewAnnotation(request.Title, request.Text, request.Tags, nil, nil),
				[]string{annotations.StorageGrafana})
		}))
	}

	if c.application.HasComponent(messengers.ComponentName) {
		<-c.application.ReadyComponent(messengers.ComponentName)
		cmp := c.application.GetComponent(messengers.ComponentName).(messengers.Component)

		subscribers = append(subscribers, mqtt.NewSubscriber(MQTTSubscribeTopicMessenger.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !c.config.Bool(boggart.ConfigMQTTMessengersEnabled) {
				return nil
			}

			parts := mqtt.RouteSplit(message.Topic())
			if len(parts) < 3 {
				return errors.New("bad topic name")
			}

			messenger := cmp.Messenger(parts[len(parts)-2])
			if messenger == nil {
				return errors.New("messenger " + parts[len(parts)-2] + " not found")
			}

			return messenger.SendMessage(parts[len(parts)-1], string(message.Payload()))
		}))
	}

	return subscribers
}

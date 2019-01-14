package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/mmcloughlin/geohash"
)

const (
	MQTTSubscribeTopicOwnTracks          mqtt.Topic = "owntracks/+/+"
	MQTTPublishTopicOwnTracksGeoHash     mqtt.Topic = "owntracks/+/+/geohash"
	MQTTSubscribeTopicAnnotationGrafana  mqtt.Topic = "annotation/grafana"
	MQTTSubscribeTopicMessenger          mqtt.Topic = "messenger/+/+"
	MQTTSubscribeTopicWOL                mqtt.Topic = boggart.ComponentName + "/wol/+"
	MQTTSubscribeTopicWOLWithIPAndSubnet mqtt.Topic = boggart.ComponentName + "/wol/+/+/+"
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
		mqtt.NewSubscriber(MQTTSubscribeTopicWOL.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !c.config.Bool(boggart.ConfigMQTTWOLEnabled) {
				return nil
			}

			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 1 {
				return errors.New("bad topic name")
			}

			mac, err := net.ParseMAC(route[len(route)-1])
			if err != nil {
				return err
			}

			return c.WOL(mac, nil, nil)
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicWOLWithIPAndSubnet.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			mac, err := net.ParseMAC(route[len(route)-3])
			if err != nil {
				return err
			}

			subnet := net.ParseIP(route[len(route)-1])
			ip := net.ParseIP(route[len(route)-2])

			return c.WOL(mac, ip, subnet)
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

func (c *Component) WOL(mac net.HardwareAddr, ip net.IP, subnet net.IP) error {
	if mac == nil {
		return errors.New("MAC isn't set")
	}

	var broadcastAddress net.IP

	if ip != nil && subnet != nil {
		broadcastAddress = net.IP{0, 0, 0, 0}
		for i := 0; i < 4; i++ {
			broadcastAddress[i] = (ip[i] & subnet[i]) | ^subnet[i]
		}
	} else {
		broadcastAddress = net.IP{255, 255, 255, 255}
	}

	c.logger.Debug("Send WOL magic packet", "mac", mac.String(), "broadcast", broadcastAddress.String())

	return wol.MagicWake(mac.String(), broadcastAddress.String())
}

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
	"github.com/mmcloughlin/geohash"
)

const (
	MQTTTopicGrafanaAnnotationGrafana mqtt.Topic = "annotation/grafana"
	MQTTTopicOwnTracks                mqtt.Topic = "owntracks/+/+"
	MQTTTopicOwnTracksGeoHash         mqtt.Topic = "owntracks/+/+/geohash"
	MQTTTopicWOL                      mqtt.Topic = boggart.ComponentName + "/wol/+"
	MQTTTopicWOLWithIPAndSubnet       mqtt.Topic = boggart.ComponentName + "/wol/+/+/+"
)

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	<-c.application.ReadyComponent(c.Name())

	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicOwnTracks.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if !c.config.Bool(boggart.ConfigOwnTracksEnabled) {
				return
			}

			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 2 {
				return
			}

			var payload map[string]interface{}

			err := json.Unmarshal(message.Payload(), &payload)
			if err != nil {
				c.logger.Warn("Failed parse request of OwnTracks", "error", err.Error())
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
			client.Publish(ctx, MQTTTopicOwnTracksGeoHash.Format(route[len(route)-2], route[len(route)-1]), message.Qos(), message.Retained(), hash)
		}),
		mqtt.NewSubscriber(MQTTTopicWOL.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !c.config.Bool(boggart.ConfigWOLEnabled) {
				return
			}

			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 1 {
				return
			}

			mac, err := net.ParseMAC(route[len(route)-1])
			if err != nil {
				c.logger.Warn("Failed parse MAC address for WOL", "error", err.Error())
				return
			}

			if err := c.WOL(mac, nil, nil); err != nil {
				c.logger.Warn("Failed send WOL magic packet", "error", err.Error())
			}
		}),
		mqtt.NewSubscriber(MQTTTopicWOLWithIPAndSubnet.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			mac, err := net.ParseMAC(route[len(route)-3])
			if err != nil {
				c.logger.Warn("Failed parse MAC address for WOL", "error", err.Error())
				return
			}

			subnet := net.ParseIP(route[len(route)-1])
			ip := net.ParseIP(route[len(route)-2])

			if err := c.WOL(mac, ip, subnet); err != nil {
				c.logger.Warn("Failed send WOL magic packet", "error", err.Error())
			}
		}),
	}

	if c.application.HasComponent(annotations.ComponentName) {
		<-c.application.ReadyComponent(annotations.ComponentName)
		cmp := c.application.GetComponent(annotations.ComponentName).(annotations.Component)

		subscribers = append(subscribers, mqtt.NewSubscriber(MQTTTopicGrafanaAnnotationGrafana.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !c.config.Bool(boggart.ConfigAnnotationsEnabled) {
				return
			}

			var request struct {
				Title string   `json:"title,omitempty"`
				Text  string   `json:"text,omitempty"`
				Tags  []string `json:"tags,omitempty"`
			}

			err := json.Unmarshal(message.Payload(), &request)
			if err != nil {
				c.logger.Warn("Failed parse request of Grafana annotation", "error", err.Error())
				return
			}

			err = cmp.CreateInStorages(
				annotations.NewAnnotation(request.Title, request.Text, request.Tags, nil, nil),
				[]string{annotations.StorageGrafana})

			if err != nil {
				c.logger.Warn("Failed send Grafana annotation", "error", err.Error())
			}
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

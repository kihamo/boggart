package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
)

const (
	MQTTTopicSimpleText        mqtt.Topic = voice.ComponentName + "/speech/+/text"
	MQTTTopicJSONText          mqtt.Topic = voice.ComponentName + "/speech/+/json"
	MQTTTopicPlayerURL         mqtt.Topic = voice.ComponentName + "/player/+/url"
	MQTTTopicPlayerAction      mqtt.Topic = voice.ComponentName + "/player/+/action"
	MQTTTopicPlayerVolume      mqtt.Topic = voice.ComponentName + "/player/+/volume"
	MQTTTopicPlayerMute        mqtt.Topic = voice.ComponentName + "/player/+/mute"
	MQTTTopicPlayerStateStatus mqtt.Topic = voice.ComponentName + "/player/+/state/status"
	MQTTTopicPlayerStateVolume mqtt.Topic = voice.ComponentName + "/player/+/state/volume"
	MQTTTopicPlayerStateMute   mqtt.Topic = voice.ComponentName + "/player/+/state/mute"
)

type SpeechRequest struct {
	Text    string  `json:"text"`
	Volume  int64   `json:"volume"`
	Speed   float64 `json:"speed"`
	Speaker string  `json:"speaker"`
}

func (c *Component) MQTTSubscribers() []mqtt.Subscriber {
	<-c.application.ReadyComponent(c.Name())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicSimpleText.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			c.Speech(ctx, route[len(route)-2], string(message.Payload()))
		}),
		mqtt.NewSubscriber(MQTTTopicJSONText.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			var request SpeechRequest

			if err := json.Unmarshal(message.Payload(), &request); err == nil {
				c.SpeechWithOptions(ctx, route[len(route)-2], request.Text, request.Volume, request.Speed, request.Speaker)
			}
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerURL.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			c.PlayURL(ctx, route[len(route)-2], string(message.Payload()))
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerAction.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			switch strings.ToLower(string(message.Payload())) {
			case "stop":
				c.Stop(ctx, route[len(route)-2])

			case "pause":
				c.Pause(ctx, route[len(route)-2])

			case "play":
				c.Play(ctx, route[len(route)-2])
			}
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerVolume.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err == nil {
				c.SetVolume(ctx, route[len(route)-2], volume)
			}
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerMute.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return
			}

			c.SetMute(ctx, route[len(route)-2], bytes.Equal(message.Payload(), []byte(`1`)))
		}),
	}
}

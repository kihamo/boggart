package internal

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
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
		mqtt.NewSubscriber(voice.MQTTTopicSimpleText, 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			c.Speech(ctx, string(message.Payload()))
		}),
		mqtt.NewSubscriber(voice.MQTTTopicJSONText, 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			var request SpeechRequest

			if err := json.Unmarshal(message.Payload(), &request); err == nil {
				c.SpeechWithOptions(ctx, request.Text, request.Volume, request.Speed, request.Speaker)
			}
		}),
		mqtt.NewSubscriber(voice.MQTTTopicPlayerURL, 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			c.PlayURL(ctx, string(message.Payload()))
		}),
		mqtt.NewSubscriber(voice.MQTTTopicPlayerAction, 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			switch strings.ToLower(string(message.Payload())) {
			case "stop":
				c.Stop(ctx)

			case "pause":
				c.Pause(ctx)

			case "play":
				c.Play(ctx)
			}
		}),
		mqtt.NewSubscriber(voice.MQTTTopicPlayerVolumeSet, 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err == nil {
				c.SetVolume(ctx, volume)
			}
		}),
	}
}

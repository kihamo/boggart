package internal

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
)

const (
	MQTTTopicSimpleText      mqtt.Topic = voice.ComponentName + "/speech/text"
	MQTTTopicJSONText        mqtt.Topic = voice.ComponentName + "/speech/json"
	MQTTTopicPlayerURL       mqtt.Topic = voice.ComponentName + "/player/url"
	MQTTTopicPlayerStatus    mqtt.Topic = voice.ComponentName + "/player/status"
	MQTTTopicPlayerAction    mqtt.Topic = voice.ComponentName + "/player/action"
	MQTTTopicPlayerVolume    mqtt.Topic = voice.ComponentName + "/player/volume"
	MQTTTopicPlayerVolumeSet mqtt.Topic = voice.ComponentName + "/player/volume/value"
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
			c.Speech(ctx, string(message.Payload()))
		}),
		mqtt.NewSubscriber(MQTTTopicJSONText.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			var request SpeechRequest

			if err := json.Unmarshal(message.Payload(), &request); err == nil {
				c.SpeechWithOptions(ctx, request.Text, request.Volume, request.Speed, request.Speaker)
			}
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerURL.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			c.PlayURL(ctx, string(message.Payload()))
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerAction.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			switch strings.ToLower(string(message.Payload())) {
			case "stop":
				c.Stop(ctx)

			case "pause":
				c.Pause(ctx)

			case "play":
				c.Play(ctx)
			}
		}),
		mqtt.NewSubscriber(MQTTTopicPlayerVolumeSet.String(), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err == nil {
				c.SetVolume(ctx, volume)
			}
		}),
	}
}

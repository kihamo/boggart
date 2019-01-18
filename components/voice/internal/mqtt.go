package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice"
)

const (
	MQTTSubscribeTopicSimpleText   mqtt.Topic = voice.ComponentName + "/speech/+/text"
	MQTTSubscribeTopicJSONText     mqtt.Topic = voice.ComponentName + "/speech/+/json"
	MQTTSubscribeTopicPlayerURL    mqtt.Topic = voice.ComponentName + "/player/+/url"
	MQTTSubscribeTopicPlayerAction mqtt.Topic = voice.ComponentName + "/player/+/action"
	MQTTSubscribeTopicPlayerVolume mqtt.Topic = voice.ComponentName + "/player/+/volume"
	MQTTSubscribeTopicPlayerMute   mqtt.Topic = voice.ComponentName + "/player/+/mute"
	MQTTTopicPlayerStateStatus     mqtt.Topic = voice.ComponentName + "/player/+/state/status"
	MQTTTopicPlayerStateVolume     mqtt.Topic = voice.ComponentName + "/player/+/state/volume"
	MQTTTopicPlayerStateMute       mqtt.Topic = voice.ComponentName + "/player/+/state/mute"
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
		mqtt.NewSubscriber(MQTTSubscribeTopicSimpleText.String(), 0, c.wrapMQTTSubscribe(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			return c.Speech(ctx, route[len(route)-2], string(message.Payload()))
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicJSONText.String(), 0, c.wrapMQTTSubscribe(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			var request SpeechRequest

			if err := json.Unmarshal(message.Payload(), &request); err != nil {
				return err
			}

			return c.SpeechWithOptions(ctx, route[len(route)-2], request.Text, request.Volume, request.Speed, request.Speaker)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlayerURL.String(), 0, c.wrapMQTTSubscribe(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			return c.PlayURL(ctx, route[len(route)-2], string(message.Payload()))
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlayerAction.String(), 0, c.wrapMQTTSubscribe(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			switch strings.ToLower(string(message.Payload())) {
			case "stop":
				return c.Stop(ctx, route[len(route)-2])

			case "pause":
				return c.Pause(ctx, route[len(route)-2])

			case "play":
				return c.Play(ctx, route[len(route)-2])
			}

			return errors.New("bad action")
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlayerVolume.String(), 0, c.wrapMQTTSubscribe(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return c.SetVolume(ctx, route[len(route)-2], volume)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlayerMute.String(), 0, c.wrapMQTTSubscribe(func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			return c.SetMute(ctx, route[len(route)-2], bytes.Equal(message.Payload(), []byte(`1`)))
		})),
	}
}

func (c *Component) wrapMQTTSubscribe(f func(ctx context.Context, client mqtt.Component, message mqtt.Message) error) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if len(c.players) == 0 {
			return nil
		}

		return f(ctx, client, message)
	}
}

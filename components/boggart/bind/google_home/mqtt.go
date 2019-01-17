package google_home

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicVolume    mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/volume"
	MQTTSubscribeTopicMute      mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/mute"
	MQTTSubscribeTopicPause     mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/pause"
	MQTTSubscribeTopicStop      mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/stop"
	MQTTSubscribeTopicPlay      mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/play"
	MQTTSubscribeTopicAction    mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/action"
	MQTTPublishTopicStateStatus mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/status"
	MQTTPublishTopicStateVolume mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/volume"
	MQTTPublishTopicStateMute   mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/mute"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicStateStatus,
		MQTTPublishTopicStateVolume,
		MQTTPublishTopicStateMute,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return b.ClientChromeCast().SetVolume(volume)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.ClientChromeCast().SetMute(bytes.Equal(message.Payload(), []byte(`1`)))
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPause.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.ClientChromeCast().Pause()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicStop.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.ClientChromeCast().Stop()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlay.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			if u := string(message.Payload()); u != "" {
				return b.ClientChromeCast().PlayFromURL(u)
			}

			return b.ClientChromeCast().Play()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicAction.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			action := string(message.Payload())

			switch strings.ToLower(action) {
			case "stop":
				return b.ClientChromeCast().Stop()

			case "pause":
				return b.ClientChromeCast().Pause()

			case "play":
				return b.ClientChromeCast().Play()
			}

			return errors.New("unknown action " + action)
		})),
	}
}

package alsa

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
	MQTTSubscribeTopicVolume    mqtt.Topic = boggart.ComponentName + "/alsa/+/volume"
	MQTTSubscribeTopicMute      mqtt.Topic = boggart.ComponentName + "/alsa/+/mute"
	MQTTSubscribeTopicPause     mqtt.Topic = boggart.ComponentName + "/alsa/+/pause"
	MQTTSubscribeTopicStop      mqtt.Topic = boggart.ComponentName + "/alsa/+/stop"
	MQTTSubscribeTopicPlay      mqtt.Topic = boggart.ComponentName + "/alsa/+/play"
	MQTTSubscribeTopicAction    mqtt.Topic = boggart.ComponentName + "/alsa/+/action"
	MQTTPublishTopicStateStatus mqtt.Topic = boggart.ComponentName + "/alsa/+/state/status"
	MQTTPublishTopicStateVolume mqtt.Topic = boggart.ComponentName + "/alsa/+/state/volume"
	MQTTPublishTopicStateMute   mqtt.Topic = boggart.ComponentName + "/alsa/+/state/mute"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicStateStatus.Format(sn)),
		mqtt.Topic(MQTTPublishTopicStateVolume.Format(sn)),
		mqtt.Topic(MQTTPublishTopicStateMute.Format(sn)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicVolume.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return b.player.SetVolume(volume)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMute.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.player.SetMute(bytes.Equal(message.Payload(), []byte(`1`)))
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPause.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.player.Pause()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicStop.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.player.Stop()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPlay.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := string(message.Payload()); u != "" {
				return b.player.PlayFromURL(u)
			}

			return b.player.Play()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicAction.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			action := string(message.Payload())

			switch strings.ToLower(action) {
			case "stop":
				return b.player.Stop()

			case "pause":
				return b.player.Pause()

			case "play":
				return b.player.Play()
			}

			return errors.New("unknown action " + action)
		})),
	}
}

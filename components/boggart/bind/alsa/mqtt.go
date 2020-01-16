package alsa

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicStateStatus,
		b.config.TopicStateVolume,
		b.config.TopicStateMute,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicVolume, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			volume, err := strconv.ParseInt(message.String(), 10, 64)
			if err != nil {
				return err
			}

			return b.SetVolume(volume)
		})),
		mqtt.NewSubscriber(b.config.TopicMute, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.SetMute(message.IsTrue())
		})),
		mqtt.NewSubscriber(b.config.TopicPause, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Pause()
		})),
		mqtt.NewSubscriber(b.config.TopicStop, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Stop()
		})),
		mqtt.NewSubscriber(b.config.TopicPlay, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := message.String(); u != "" {
				return b.PlayFromURL(u)
			}

			return b.Play()
		})),
		mqtt.NewSubscriber(b.config.TopicResume, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := message.String(); u != "" {
				return b.PlayFromURL(u)
			}

			return b.Play()
		})),
		mqtt.NewSubscriber(b.config.TopicAction, 0, b.MQTTContainer().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			action := message.String()

			switch strings.ToLower(action) {
			case "stop":
				return b.Stop()

			case "pause":
				return b.Pause()

			case "play", "resume":
				return b.Play()
			}

			return errors.New("unknown action " + action)
		})),
	}
}

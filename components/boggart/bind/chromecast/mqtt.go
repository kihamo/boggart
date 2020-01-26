package chromecast

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
		b.config.TopicStateContent,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicVolume, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			volume, err := strconv.ParseInt(message.String(), 10, 64)
			if err != nil {
				return err
			}

			return b.SetVolume(ctx, volume)
		})),
		mqtt.NewSubscriber(b.config.TopicMute, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.SetMute(ctx, message.IsTrue())
		})),
		mqtt.NewSubscriber(b.config.TopicPause, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Pause(ctx)
		})),
		mqtt.NewSubscriber(b.config.TopicStop, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Stop(ctx)
		})),
		mqtt.NewSubscriber(b.config.TopicPlay, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := message.String(); u != "" {
				return b.PlayFromURL(ctx, u)
			}

			return b.Resume(ctx)
		})),
		mqtt.NewSubscriber(b.config.TopicSeek, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			second, err := strconv.ParseUint(message.String(), 10, 64)
			if err != nil {
				return err
			}

			return b.Seek(ctx, second)
		})),
		mqtt.NewSubscriber(b.config.TopicResume, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Resume(ctx)
		})),
		mqtt.NewSubscriber(b.config.TopicAction, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			action := message.String()

			switch strings.ToLower(action) {
			case "stop":
				return b.Stop(ctx)

			case "pause":
				return b.Pause(ctx)

			case "play", "resume":
				return b.Resume(ctx)
			}

			return errors.New("unknown action " + action)
		})),
	}
}

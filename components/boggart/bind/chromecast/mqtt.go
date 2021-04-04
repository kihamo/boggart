package chromecast

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	id := b.Meta().ID()
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicVolume.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			volume, err := strconv.ParseInt(message.String(), 10, 64)
			if err != nil {
				return err
			}

			return b.SetVolume(ctx, volume)
		})),
		mqtt.NewSubscriber(cfg.TopicMute.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.SetMute(ctx, message.IsTrue())
		})),
		mqtt.NewSubscriber(cfg.TopicPause.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Pause(ctx)
		})),
		mqtt.NewSubscriber(cfg.TopicStop.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Stop(ctx)
		})),
		mqtt.NewSubscriber(cfg.TopicPlay.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := message.String(); u != "" {
				return b.PlayFromURL(ctx, u)
			}

			return b.Resume(ctx)
		})),
		mqtt.NewSubscriber(cfg.TopicSeek.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			second, err := strconv.ParseUint(message.String(), 10, 64)
			if err != nil {
				return err
			}

			return b.Seek(ctx, second)
		})),
		mqtt.NewSubscriber(cfg.TopicResume.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Resume(ctx)
		})),
		mqtt.NewSubscriber(cfg.TopicAction.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
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

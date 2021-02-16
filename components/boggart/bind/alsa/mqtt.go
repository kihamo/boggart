package alsa

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	sn := b.Meta().SerialNumber()
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicVolume.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			volume, err := strconv.ParseInt(message.String(), 10, 64)
			if err != nil {
				return err
			}

			return b.SetVolume(volume)
		})),
		mqtt.NewSubscriber(cfg.TopicMute.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.SetMute(message.IsTrue())
		})),
		mqtt.NewSubscriber(cfg.TopicPause.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Pause()
		})),
		mqtt.NewSubscriber(cfg.TopicStop.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Stop()
		})),
		mqtt.NewSubscriber(cfg.TopicPlay.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := message.String(); u != "" {
				return b.PlayFromURL(u)
			}

			return b.Play()
		})),
		mqtt.NewSubscriber(cfg.TopicResume.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if u := message.String(); u != "" {
				return b.PlayFromURL(u)
			}

			return b.Play()
		})),
		mqtt.NewSubscriber(cfg.TopicAction.Format(sn), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
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

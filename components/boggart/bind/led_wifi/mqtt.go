package ledwifi

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/wifiled"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicPower.Format(cfg.Address), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if message.IsTrue() {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
		mqtt.NewSubscriber(cfg.TopicColor.Format(cfg.Address), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			color, err := wifiled.ColorFromString(message.String())
			if err != nil {
				return err
			}

			return b.bulb.SetColorPersist(ctx, *color)
		})),
		mqtt.NewSubscriber(cfg.TopicMode.Format(cfg.Address), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			mode, err := wifiled.ModeFromString(strings.TrimSpace(message.String()))
			if err != nil {
				return err
			}

			state, err := b.bulb.State(ctx)
			if err != nil {
				return err
			}

			return b.bulb.SetMode(ctx, *mode, state.Speed)
		})),
		mqtt.NewSubscriber(cfg.TopicSpeed.Format(cfg.Address), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			speed, err := strconv.ParseInt(strings.TrimSpace(message.String()), 10, 64)
			if err != nil {
				return err
			}

			state, err := b.bulb.State(ctx)
			if err != nil {
				return err
			}

			return b.bulb.SetMode(ctx, state.Mode, uint8(speed))
		})),
	}
}

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
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicPower.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
			if message.IsTrue() {
				err = b.bulb.PowerOn(ctx)
			} else {
				err = b.bulb.PowerOff(ctx)
			}

			if err == nil {
				err = b.runUpdater(ctx)
			}

			return err
		})),
		mqtt.NewSubscriber(cfg.TopicColor.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			color, err := wifiled.ColorFromString(message.String())
			if err != nil {
				return err
			}

			err = b.bulb.SetColorPersist(ctx, *color)
			if err == nil {
				err = b.runUpdater(ctx)
			}

			return err
		})),
		mqtt.NewSubscriber(cfg.TopicMode.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			mode, err := wifiled.ModeFromString(strings.TrimSpace(message.String()))
			if err != nil {
				return err
			}

			state, err := b.bulb.State(ctx)
			if err != nil {
				return err
			}

			err = b.bulb.SetMode(ctx, *mode, state.Speed)
			if err == nil {
				err = b.runUpdater(ctx)
			}

			return err
		})),
		mqtt.NewSubscriber(cfg.TopicSpeed.Format(id), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			speed, err := strconv.ParseInt(strings.TrimSpace(message.String()), 10, 64)
			if err != nil {
				return err
			}

			state, err := b.bulb.State(ctx)
			if err != nil {
				return err
			}

			err = b.bulb.SetMode(ctx, state.Mode, uint8(speed))
			if err == nil {
				err = b.runUpdater(ctx)
			}

			return err
		})),
	}
}

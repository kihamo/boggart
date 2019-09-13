package led_wifi

import (
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/wifiled"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/led/+/"

	MQTTSubscribeTopicPower       = MQTTPrefix + "power"
	MQTTSubscribeTopicColor       = MQTTPrefix + "color"
	MQTTSubscribeTopicMode        = MQTTPrefix + "mode"
	MQTTSubscribeTopicSpeed       = MQTTPrefix + "speed"
	MQTTPublishTopicStatePower    = MQTTPrefix + "state/power"
	MQTTPublishTopicStateColor    = MQTTPrefix + "state/color"
	MQTTPublishTopicStateColorHSV = MQTTPrefix + "state/color/hsv"
	MQTTPublishTopicStateMode     = MQTTPrefix + "state/mode"
	MQTTPublishTopicStateSpeed    = MQTTPrefix + "state/speed"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	host := b.bulb.Host()

	return []mqtt.Topic{
		MQTTPublishTopicStatePower.Format(host),
		MQTTPublishTopicStateColor.Format(host),
		MQTTPublishTopicStateColorHSV.Format(host),
		MQTTPublishTopicStateMode.Format(host),
		MQTTPublishTopicStateSpeed.Format(host),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	host := mqtt.NameReplace(b.bulb.Host())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicPower.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if message.IsTrue() {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicColor.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			color, err := wifiled.ColorFromString(message.String())
			if err != nil {
				return err
			}

			return b.bulb.SetColorPersist(ctx, *color)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMode.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
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
		mqtt.NewSubscriber(MQTTSubscribeTopicSpeed.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
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

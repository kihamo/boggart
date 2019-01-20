package led_wifi

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/wifiled"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicPower       mqtt.Topic = boggart.ComponentName + "/led/+/power"
	MQTTSubscribeTopicColor       mqtt.Topic = boggart.ComponentName + "/led/+/color"
	MQTTSubscribeTopicMode        mqtt.Topic = boggart.ComponentName + "/led/+/mode"
	MQTTSubscribeTopicSpeed       mqtt.Topic = boggart.ComponentName + "/led/+/speed"
	MQTTPublishTopicStatePower    mqtt.Topic = boggart.ComponentName + "/led/+/state/power"
	MQTTPublishTopicStateColor    mqtt.Topic = boggart.ComponentName + "/led/+/state/color"
	MQTTPublishTopicStateColorHSV mqtt.Topic = boggart.ComponentName + "/led/+/state/color/hsv"
	MQTTPublishTopicStateMode     mqtt.Topic = boggart.ComponentName + "/led/+/state/mode"
	MQTTPublishTopicStateSpeed    mqtt.Topic = boggart.ComponentName + "/led/+/state/speed"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	host := mqtt.NameReplace(b.bulb.Host())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicStatePower.Format(host)),
		mqtt.Topic(MQTTPublishTopicStateColor.Format(host)),
		mqtt.Topic(MQTTPublishTopicStateColorHSV.Format(host)),
		mqtt.Topic(MQTTPublishTopicStateMode.Format(host)),
		mqtt.Topic(MQTTPublishTopicStateSpeed.Format(host)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	host := mqtt.NameReplace(b.bulb.Host())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicPower.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if bytes.Equal(message.Payload(), []byte(`1`)) {
				return b.On(ctx)
			}

			return b.Off(ctx)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicColor.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			color, err := wifiled.ColorFromString(string(message.Payload()))
			if err != nil {
				return err
			}

			return b.bulb.SetColorPersist(ctx, *color)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMode.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			mode, err := wifiled.ModeFromString(strings.TrimSpace(string(message.Payload())))
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
			speed, err := strconv.ParseInt(strings.TrimSpace(string(message.Payload())), 10, 64)
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

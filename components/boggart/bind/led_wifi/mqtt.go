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
	MQTTTopicPower         mqtt.Topic = boggart.ComponentName + "/led/+/power"
	MQTTTopicColor         mqtt.Topic = boggart.ComponentName + "/led/+/color"
	MQTTTopicMode          mqtt.Topic = boggart.ComponentName + "/led/+/mode"
	MQTTTopicSpeed         mqtt.Topic = boggart.ComponentName + "/led/+/speed"
	MQTTTopicStatePower    mqtt.Topic = boggart.ComponentName + "/led/+/state/power"
	MQTTTopicStateColor    mqtt.Topic = boggart.ComponentName + "/led/+/state/color"
	MQTTTopicStateColorHSV mqtt.Topic = boggart.ComponentName + "/led/+/state/color/hsv"
	MQTTTopicStateMode     mqtt.Topic = boggart.ComponentName + "/led/+/state/mode"
	MQTTTopicStateSpeed    mqtt.Topic = boggart.ComponentName + "/led/+/state/speed"
)

func (b *Bind) MQTTTopics() []mqtt.Topic {
	host := mqtt.NameReplace(b.bulb.Host())

	return []mqtt.Topic{
		mqtt.Topic(MQTTTopicPower.Format(host)),
		mqtt.Topic(MQTTTopicColor.Format(host)),
		mqtt.Topic(MQTTTopicMode.Format(host)),
		mqtt.Topic(MQTTTopicSpeed.Format(host)),
		mqtt.Topic(MQTTTopicStatePower.Format(host)),
		mqtt.Topic(MQTTTopicStateColor.Format(host)),
		mqtt.Topic(MQTTTopicStateColorHSV.Format(host)),
		mqtt.Topic(MQTTTopicStateMode.Format(host)),
		mqtt.Topic(MQTTTopicStateSpeed.Format(host)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	host := mqtt.NameReplace(b.bulb.Host())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicPower.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
			if bytes.Equal(message.Payload(), []byte(`1`)) {
				b.On(ctx)
			} else {
				b.Off(ctx)
			}
		})),
		mqtt.NewSubscriber(MQTTTopicColor.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
			color, err := wifiled.ColorFromString(string(message.Payload()))
			if err != nil {
				return
			}

			b.bulb.SetColorPersist(ctx, *color)
		})),
		mqtt.NewSubscriber(MQTTTopicMode.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
			mode, err := wifiled.ModeFromString(strings.TrimSpace(string(message.Payload())))
			if err != nil {
				return
			}

			state, err := b.bulb.State(ctx)
			if err != nil {
				return
			}

			b.bulb.SetMode(ctx, *mode, state.Speed)
		})),
		mqtt.NewSubscriber(MQTTTopicSpeed.Format(host), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
			speed, err := strconv.ParseInt(strings.TrimSpace(string(message.Payload())), 10, 64)
			if err != nil {
				return
			}

			state, err := b.bulb.State(ctx)
			if err != nil {
				return
			}

			b.bulb.SetMode(ctx, state.Mode, uint8(speed))
		})),
	}
}

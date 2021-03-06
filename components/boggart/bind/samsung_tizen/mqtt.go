package tizen

import (
	"context"
	"errors"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicPower, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			if message.IsTrue() {
				if mac := b.Meta().MAC(); mac != nil {
					return wol.MagicWake(mac.String(), "255.255.255.255")
				}
			}

			if !b.Meta().Status().IsStatusOnline() {
				return errors.New("bind isn't online")
			}

			return b.client.SendCommand(tv.KeyPower)
		}),
		mqtt.NewSubscriber(cfg.TopicKey, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			return b.client.SendCommand(message.String())
		})),
	}
}

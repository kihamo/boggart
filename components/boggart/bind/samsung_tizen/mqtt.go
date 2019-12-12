package samsung_tizen

import (
	"context"
	"errors"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicDeviceID,
		b.config.TopicDeviceModelName,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicPower, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			if message.IsTrue() {
				b.mutex.RLock()
				mac := b.mac
				b.mutex.RUnlock()

				return wol.MagicWake(mac, "255.255.255.255")
			}

			if !b.IsStatusOnline() {
				return errors.New("bind isn't online")
			}

			return b.client.SendCommand(tv.KeyPower)
		}),
		mqtt.NewSubscriber(b.config.TopicKey, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.client.SendCommand(message.String())
		})),
	}
}

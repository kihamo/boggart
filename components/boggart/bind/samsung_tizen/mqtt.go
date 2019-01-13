package samsung_tizen

import (
	"bytes"
	"context"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicPower           mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	MQTTTopicKey             mqtt.Topic = boggart.ComponentName + "/tv/+/key"
	MQTTTopicDeviceID        mqtt.Topic = boggart.ComponentName + "/tv/+/device/id"
	MQTTTopicDeviceModelName mqtt.Topic = boggart.ComponentName + "/tv/+/device/model-name"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicPower.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				b.mutex.RLock()
				mac := b.mac
				b.mutex.RUnlock()

				wol.MagicWake(mac, "255.255.255.255")
			} else if b.Status() == boggart.BindStatusOnline {
				b.client.SendCommand(tv.KeyPower)
			}
		}),
		mqtt.NewSubscriber(MQTTTopicKey.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			b.client.SendCommand(string(message.Payload()))
		})),
	}
}

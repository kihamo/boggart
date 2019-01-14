package google_home

import (
	"bytes"
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicVolume mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/volume"
	MQTTSubscribeTopicMute   mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/mute"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return b.ClientChromeCast().SetVolume(volume)
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.ClientChromeCast().SetMute(bytes.Equal(message.Payload(), []byte(`1`)))
		})),
	}
}

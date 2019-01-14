package google_home

import (
	"bytes"
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicVolume      mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/volume"
	MQTTTopicMute        mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/mute"
	MQTTTopicStateVolume mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/volume"
	MQTTTopicStateMute   mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/mute"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTTopicVolume,
		MQTTTopicMute,
		MQTTTopicStateVolume,
		MQTTTopicStateMute,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return err
			}

			return b.ClientChromeCast().SetVolume(volume)
		})),
		mqtt.NewSubscriber(MQTTTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.ClientChromeCast().SetMute(bytes.Equal(message.Payload(), []byte(`1`)))
		})),
	}
}

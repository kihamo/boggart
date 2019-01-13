package lg_webos

import (
	"bytes"
	"context"
	"strconv"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicApplication        mqtt.Topic = boggart.ComponentName + "/tv/+/application"
	MQTTTopicMute               mqtt.Topic = boggart.ComponentName + "/tv/+/mute"
	MQTTTopicVolume             mqtt.Topic = boggart.ComponentName + "/tv/+/volume"
	MQTTTopicVolumeUp           mqtt.Topic = boggart.ComponentName + "/tv/+/volume/up"
	MQTTTopicVolumeDown         mqtt.Topic = boggart.ComponentName + "/tv/+/volume/down"
	MQTTTopicToast              mqtt.Topic = boggart.ComponentName + "/tv/+/toast"
	MQTTTopicPower              mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	MQTTTopicStateMute          mqtt.Topic = boggart.ComponentName + "/tv/+/state/mute"
	MQTTTopicStateVolume        mqtt.Topic = boggart.ComponentName + "/tv/+/state/volume"
	MQTTTopicStateApplication   mqtt.Topic = boggart.ComponentName + "/tv/+/state/application"
	MQTTTopicStateChannelNumber mqtt.Topic = boggart.ComponentName + "/tv/+/state/channel-number"
	MQTTTopicStatePower         mqtt.Topic = boggart.ComponentName + "/tv/+/state/power"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTTopicStateMute,
		MQTTTopicStateVolume,
		MQTTTopicStateApplication,
		MQTTTopicStateChannelNumber,
		MQTTTopicStatePower,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTTopicApplication.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			if client, err := b.Client(); err == nil {
				client.ApplicationManagerLaunch(string(message.Payload()), nil)
			}
		})),
		mqtt.NewSubscriber(MQTTTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			if client, err := b.Client(); err == nil {
				client.AudioSetMute(bytes.Equal(message.Payload(), []byte(`1`)))
			}
		})),
		mqtt.NewSubscriber(MQTTTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			vol, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err == nil {
				if client, err := b.Client(); err == nil {
					client.AudioSetVolume(int(vol))
				}
			}
		})),
		mqtt.NewSubscriber(MQTTTopicVolumeUp.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
				return
			}

			if client, err := b.Client(); err == nil {
				client.AudioVolumeUp()
			}
		})),
		mqtt.NewSubscriber(MQTTTopicVolumeDown.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
				return
			}

			if client, err := b.Client(); err == nil {
				client.AudioVolumeDown()
			}
		})),
		mqtt.NewSubscriber(MQTTTopicToast.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			if client, err := b.Client(); err == nil {
				client.SystemNotificationsCreateToast(string(message.Payload()))
			}
		})),
		mqtt.NewSubscriber(MQTTTopicPower.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return
			}

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				wol.MagicWake(b.SerialNumber(), "255.255.255.255")
			} else if b.Status() == boggart.BindStatusOnline {
				if client, err := b.Client(); err == nil {
					client.SystemTurnOff()
				}
			}
		}),
	}
}

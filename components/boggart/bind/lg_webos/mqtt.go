package lg_webos

import (
	"context"
	"errors"
	"strconv"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicApplication      mqtt.Topic = boggart.ComponentName + "/tv/+/application"
	MQTTSubscribeTopicMute             mqtt.Topic = boggart.ComponentName + "/tv/+/mute"
	MQTTSubscribeTopicVolume           mqtt.Topic = boggart.ComponentName + "/tv/+/volume"
	MQTTSubscribeTopicVolumeUp         mqtt.Topic = boggart.ComponentName + "/tv/+/volume/up"
	MQTTSubscribeTopicVolumeDown       mqtt.Topic = boggart.ComponentName + "/tv/+/volume/down"
	MQTTSubscribeTopicToast            mqtt.Topic = boggart.ComponentName + "/tv/+/toast"
	MQTTSubscribeTopicPower            mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	MQTTPublishTopicStateMute          mqtt.Topic = boggart.ComponentName + "/tv/+/state/mute"
	MQTTPublishTopicStateVolume        mqtt.Topic = boggart.ComponentName + "/tv/+/state/volume"
	MQTTPublishTopicStateApplication   mqtt.Topic = boggart.ComponentName + "/tv/+/state/application"
	MQTTPublishTopicStateChannelNumber mqtt.Topic = boggart.ComponentName + "/tv/+/state/channel-number"
	MQTTPublishTopicStatePower         mqtt.Topic = boggart.ComponentName + "/tv/+/state/power"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicStateMute,
		MQTTPublishTopicStateVolume,
		MQTTPublishTopicStateApplication,
		MQTTPublishTopicStateChannelNumber,
		MQTTPublishTopicStatePower,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicApplication, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			client, err := b.Client()
			if err != nil {
				return err
			}

			_, err = client.ApplicationManagerLaunch(message.String(), nil)
			return err
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicMute, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			client, err := b.Client()
			if err != nil {
				return err
			}

			return client.AudioSetMute(message.IsTrue())
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicVolume, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			vol, err := strconv.ParseInt(message.String(), 10, 64)
			if err != nil {
				return err
			}

			client, err := b.Client()
			if err != nil {
				return err
			}

			return client.AudioSetVolume(int(vol))
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicVolumeUp, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
				return nil
			}

			client, err := b.Client()
			if err != nil {
				return err
			}

			return client.AudioVolumeUp()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicVolumeDown, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
				return nil
			}

			client, err := b.Client()
			if err != nil {
				return err
			}

			return client.AudioVolumeDown()
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicToast, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			return b.Toast(message.String())
		})),
		mqtt.NewSubscriber(MQTTSubscribeTopicPower, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 2) {
				return nil
			}

			if message.IsTrue() {
				return wol.MagicWake(b.SerialNumber(), "255.255.255.255")
			}

			if b.Status() != boggart.BindStatusOnline {
				return errors.New("bind isn't online")
			}

			client, err := b.Client()
			if err != nil {
				return err
			}

			return client.SystemTurnOff()
		}),
	}
}

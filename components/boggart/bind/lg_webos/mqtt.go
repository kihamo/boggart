package webos

import (
	"context"
	"errors"
	"strconv"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	mac := b.Meta().MACAsString()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicApplication.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			client := b.Client()
			if client == nil {
				return errors.New("client isn't init")
			}

			_, err := client.ApplicationManagerLaunch(message.String(), nil)
			return err
		})),
		mqtt.NewSubscriber(cfg.TopicMute.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			client := b.Client()
			if client == nil {
				return errors.New("client isn't init")
			}

			return client.AudioSetMute(message.IsTrue())
		})),
		mqtt.NewSubscriber(cfg.TopicVolume.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			vol, err := strconv.ParseInt(message.String(), 10, 64)
			if err != nil {
				return err
			}

			client := b.Client()
			if client == nil {
				return errors.New("client isn't init")
			}

			return client.AudioSetVolume(int(vol))
		})),
		mqtt.NewSubscriber(cfg.TopicVolumeUp.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -3) {
				return nil
			}

			client := b.Client()
			if client == nil {
				return errors.New("client isn't init")
			}

			return client.AudioVolumeUp()
		})),
		mqtt.NewSubscriber(cfg.TopicVolumeDown.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -3) {
				return nil
			}

			client := b.Client()
			if client == nil {
				return errors.New("client isn't init")
			}

			return client.AudioVolumeDown()
		})),
		mqtt.NewSubscriber(cfg.TopicToast.Format(mac), 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			return b.Toast(message.String())
		})),
		mqtt.NewSubscriber(cfg.TopicPower.Format(mac), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckMACInTopic(message.Topic(), -2) {
				return nil
			}

			if message.IsTrue() {
				return wol.MagicWake(b.Meta().MACAsString(), "255.255.255.255")
			}

			if !b.Meta().Status().IsStatusOnline() {
				return errors.New("bind isn't online")
			}

			client := b.Client()
			if client == nil {
				return errors.New("client isn't init")
			}

			return client.SystemTurnOff()
		}),
	}
}

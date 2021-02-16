package miio

import (
	"context"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"go.uber.org/multierr"
)

var (
	payloadDone = []byte(`done`)
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicSetFanPower, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSetFanPower)),
		mqtt.NewSubscriber(cfg.TopicSetVolume, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSetVolume)),
		mqtt.NewSubscriber(cfg.TopicTestVolume, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTTestVolume)),
		mqtt.NewSubscriber(cfg.TopicFind, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTFind)),
		mqtt.NewSubscriber(cfg.TopicAction, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTAction)),
	}
}

func (b *Bind) callbackMQTTSetFanPower(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -3) {
		return nil
	}

	power, err := strconv.ParseUint(message.String(), 10, 64)
	if err != nil {
		return err
	}

	err = b.device.SetFanPower(ctx, power)
	if err == nil {
		err = b.updateFanPower(ctx)
	}

	return err
}

func (b *Bind) callbackMQTTSetVolume(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -3) {
		return nil
	}

	power, err := strconv.ParseUint(message.String(), 10, 64)
	if err != nil {
		return err
	}

	err = b.device.SetSoundVolume(ctx, uint32(power))
	if err == nil {
		err = b.updateVolume(ctx)
	}

	return err
}

func (b *Bind) callbackMQTTTestVolume(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !message.IsTrue() || !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -3) {
		return nil
	}

	err := b.device.SoundVolumeTest(ctx)
	if err != nil {
		return err
	}

	return b.MQTT().PublishAsyncRawWithoutCache(ctx, message.Topic(), 1, false, payloadDone)
}

func (b *Bind) callbackMQTTFind(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !message.IsTrue() || !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -2) {
		return nil
	}

	err := b.device.Find(ctx)
	if err != nil {
		return err
	}

	return b.MQTT().PublishAsyncRawWithoutCache(ctx, message.Topic(), 1, false, payloadDone)
}

func (b *Bind) callbackMQTTAction(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -2) {
		return nil
	}

	var err error

	switch message.String() {
	case "start":
		err = b.device.Start(ctx)
	case "spot":
		err = b.device.Spot(ctx)
	case "stop":
		err = b.device.Stop(ctx)
	case "pause":
		err = b.device.Pause(ctx)
	case "home":
		err = b.device.Home(ctx)
	}

	if err == nil {
		time.Sleep(time.Second * 5)

		err = b.updateState(ctx)
	}

	return err
}

func (b *Bind) updateState(ctx context.Context) error {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	status, err := b.device.Status(ctx)
	if err == nil {
		cfg := b.config()

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicState.Format(sn), status.State); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicError.Format(sn), status.Error); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

func (b *Bind) updateFanPower(ctx context.Context) error {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	fan, err := b.device.FanPower(ctx)
	if err == nil {
		err = b.MQTT().PublishAsync(ctx, b.config().TopicFanPower.Format(sn), fan)
	}

	return err
}

func (b *Bind) updateVolume(ctx context.Context) error {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	volume, err := b.device.SoundVolume(ctx)
	if err == nil {
		err = b.MQTT().PublishAsync(ctx, b.config().TopicVolume.Format(sn), volume)
	}

	return err
}

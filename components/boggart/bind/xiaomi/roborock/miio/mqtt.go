package miio

import (
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/mqtt"
	"go.uber.org/multierr"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicBattery,
		b.config.TopicCleanArea,
		b.config.TopicCleanTime,
		b.config.TopicFanPower,
		b.config.TopicConsumableFilter,
		b.config.TopicConsumableBrushMain,
		b.config.TopicConsumableBrushSide,
		b.config.TopicConsumableSensor,
		b.config.TopicState,
		b.config.TopicError,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicSetFanPower, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSetFanPower)),
		mqtt.NewSubscriber(b.config.TopicSetVolume, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSetVolume)),
		mqtt.NewSubscriber(b.config.TopicTestVolume, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTTestVolume)),
		mqtt.NewSubscriber(b.config.TopicFind, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTFind)),
		mqtt.NewSubscriber(b.config.TopicAction, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTAction)),
	}
}

func (b *Bind) callbackMQTTSetFanPower(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 3) {
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
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 3) {
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
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 3) {
		return nil
	}

	return b.device.SoundVolumeTest(ctx)
}

func (b *Bind) callbackMQTTFind(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 2) {
		return nil
	}

	return b.device.Find(ctx)
}

func (b *Bind) callbackMQTTAction(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 2) {
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
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicState.Format(sn), status.State); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicError.Format(sn), status.Error); e != nil {
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
		err = b.MQTT().PublishAsync(ctx, b.config.TopicFanPower.Format(sn), fan)
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
		err = b.MQTT().PublishAsync(ctx, b.config.TopicVolume.Format(sn), volume)
	}

	return err
}
